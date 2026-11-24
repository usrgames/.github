# testdata/retry

Runs saved from GitHub's API for `ci.yml` to rehearse `retry.yml` on, trimmed to the fields it reads: the run (`run.json`), the newest run of its workflow on its branch when it ended (`newest.json`), its jobs (`jobs.json`), and the annotations of each job cancelled with no steps (`annotations-JOB.json`). A run that ended `cancelled` holds only `run.json`, since nothing else is read.

| Folder | Run | Decided |
|---|---|---|
| `no-runner-glonk-ci` | glonk 45262717999, attempt 1 | rerun |
| `no-runner-tty1-release` | tty1 45450645537, attempt 1, the release of v0.2.0 | rerun |
| `no-runner-site-ci` | usr.games 45460600664, attempt 1 | rerun |
| `no-runner-github-ci` | .github 45469298847, attempt 1 | rerun |
| `failure-github-ci` | .github 45061923679, a build that failed | leave |
| `superseded-glonk-ci` | glonk 42575882407, cancelled by a newer push | leave |
| `newer-run` | tty1's, with a newer run on its branch | leave |
| `third-attempt` | tty1's, as its third attempt | leave |
| `other-cancel` | tty1's, with its job cancelled for another reason ("The operation was canceled.") | leave |

The last three are made from `no-runner-tty1-release` by changing one field each, since no saved run shows them.
