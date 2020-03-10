package main

import (
  "errors"
  "fmt"
  "os/exec"
  "time"

  "github.com/spf13/cobra"
)

var (
  errWorkingTreeIsNotClean = errors.New("Working tree not clean")
  errMergeConflict = errors.New("Merge conflicted")
)

func main() {
  c := initContext()

  var noPush bool
  var noFix bool

  rootCmd := &cobra.Command {
    Use: "quickfollow [Repository Path]",
    Short: "Follows upstream branch",
    Args: cobra.MaximumNArgs(1),
    TraverseChildren: true,
    Run: func(cmd *cobra.Command, args []string) {
      path := "."
      if len(args) > 0 {
        path = args[0]
      }

      c.info("==QuickFollow Start==")
      act(c, path, noPush, noFix)

      c.next("E")
      c.info("==QuickFollow End==")
    },
  }
  rootCmd.Flags().BoolVar(&noPush, "no-push", false, "Don't execute \"git push --all\"")
  rootCmd.Flags().BoolVar(&noFix, "no-fix", false, "Cause panic when merge conflicts")

  rootCmd.Execute()
}

func act(c *context, path string, noPush bool, noFix bool) {
  c.next("S")
  config := load(c, path)

  c.next("F")
  fetch(c, path, config)

  c.next("M")
  requireMerge := followAll(c, config, path, noFix)

  c.next("J")
  joinAll(c, path, requireMerge)

  if noPush {
    return
  }

  c.next("P")
  push(c, path)
}

func fetch(c *context, path string, config config) {
  if config.remote != "" {
    c.info(">> git remote add %s %s", config.upstream, config.remote)
    cmd(path, "git", "remote", "add", config.upstream, config.remote).Run()
    c.info("<< done")
  }
  c.info(">> git fetch %s", config.upstream)
  cmd(path, "git", "fetch", config.upstream).Run()
  c.info("<< done")
}

func followAll(c *context, config config, path string, noFix bool) []string {
  c.info("Total target branch: %d", len(config.target))
  requireMerge := make([]string, 0)
  for i, n := range config.target {
    c.step("S")
    if !isClean(path) {
      c.panic(errWorkingTreeIsNotClean)
    }
    c.info("Next branch %s (%d/%d)", n, i + 1, len(config.target))
    if follow(c, config, path, n, noFix) {
      requireMerge = append(requireMerge, n)
    }
    c.back()
  }
  if len(requireMerge) == 0 {
    requireMerge = append(requireMerge, config.target[len(config.target) - 1])
  }
  return requireMerge
}

func isClean(path string) bool {
  out, err := cmd(path, "git", "status", "--porcelain").Output()
  return err == nil && len(out) == 0
}

func cmd(path string, command string, a ...string) *exec.Cmd {
  cmd := exec.Command(command, a...)
  cmd.Dir = path
  return cmd
}

func follow(c *context, config config, path string, name string, noFix bool) bool {
  c.next("C")
  checkout(c, path, name)

  c.next("M")
  merge(c, path, fmt.Sprintf("%s/%s", config.upstream, config.branch))
  if !isClean(path) {
    if noFix {
      c.panic(errMergeConflict)
    }
    c.info("Merge conflicted! Please resolve all conflicts and commit")
    waitCommit(path)
    c.info("Conflict resolved. Working tree clean")
    return true
  }
  return false
}

func checkout(c *context, path string, target string) {
  c.info(">> git checkout %s", target)
  err := cmd(path, "git", "checkout", target).Run()
  if err != nil {
    c.panic(err)
  }
  c.info("<< done")
}

func merge(c *context, path string, target string) {
  c.info(">> git merge %s", target)
  cmd(path, "git", "merge", target).Run()
  c.info("<< done")
}

func waitCommit(path string) {
  for !isClean(path) {
    time.Sleep(3 * time.Second)
  }
}

func joinAll(c *context, path string, requireMerge []string) {
  checkout(c, path, "master")
  c.info("Branches which requires merge: %d", len(requireMerge))
  for i, n := range requireMerge {
    if !isClean(path) {
      c.panic(errWorkingTreeIsNotClean)
    }
    c.info("Next branch %s (%d/%d)", n, i + 1, len(requireMerge))
    join(c, path, n)
  }
}

func join(c *context, path string, n string) {
  merge(c, path, n)
  if !isClean(path) {
    c.info("Merge conflicted! Please resolve all conflicts and commit")
    waitCommit(path)
    c.info("Conflict resolved. Working tree clean")
  }
}

func push(c *context, path string) {
  c.info(">> git push --all")
  cmd(path, "git", "push", "--all").Run()
  c.info("<< done")
}
