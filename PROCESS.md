# Process

This document is to describe the internal process that the Vega project core protocol team uses for issue and pull request management and the development workflow.

## Table of contents
 * [Project Boards](#project-boards)
 * [Workflow](#workflow)
 * [Reviews](#reviews)

## Project Boards

Work is primarily driven by the ~3 month [milestone planning](https://github.com/vegaprotocol/specs-internal/tree/master/milestones) and split into 2 week sprints for the team have short term delivery focus. The Core Protocol [Kanban Board](https://github.com/vegaprotocol/vega/projects?type=beta) uses two GitHub Actions to automate some of the project management:

### [Add Issues To Project Board](https://github.com/vegaprotocol/vega/tree/develop/.github/workflows/add_issue_to_project.yml) GitHub Action

Any issue that is opened in the following repos will be placed in the Sprint Backlog on the [Core Kanban Board](https://github.com/orgs/vegaprotocol/projects/106/views/27) set as `no:status` in order for it be reviewed, refined and planned into a Sprint. Issues added to the core project board will be auto labelled with `wallet`.

### [Manage Project Board](https://github.com/vegaprotocol/vega/tree/develop/.github/workflows/project_manage.yml) GitHub Action

For any pull request to be merged into develop or main/master matching a pull request action type specified in the GitHub Action the following jobs are run. The `no-changelog` label can be applied to the pull request for times when we do not need a changelog entry or a verifiable linked issue.

> ℹ️ Note: Adding the `no-changelog` label to a pull request should only be done for smaller (> 1 hour) tasks that do not require a changelog entry or a verifiable linked issue.

#### Verify Conventional Commits Job

This job will run on all pull requests that are not created by the Renovate Bot. A [third party GitHub Action](https://commitsar.aevea.ee/usage/github/) checks sure all commit messages adhere to the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0-beta.2/) specification.

#### Verify Linked Issue Job

This job will run on all pull requests that are not created by the Renovate Bot and that do not have the `no-changelog` label applied. Adding and removing this label from any pull request will have the following effect:

* **adding** the `no-changelog` label will **stop this job from running**
* **removing** the `no-changelog` label will **start running this job**

A [third party GitHub action](https://github.com/hattan/verify-linked-issue-action) makes sure all pull requests that meet the above criteria have a linked issue.

> ℹ️ Note: Linked issues must have a space after the issue number in the first comment of the PR i.e. “Closes #123 “. For further details on linking keywords see the [GitHub docs](https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue#linking-a-pull-request-to-an-issue-using-a-keyword)

#### Verify CHANGELOG Updated Job

This job will run on all **non-draft** pull requests that are not created by the Renovate Bot. Adding and removing the `no-changelog` label from a non-draft pull request will have the following effect:

* **adding** the `no-changelog` label will **stop this job from running**
* **removing** the `no-changelog` label will **start running this job**

A [third party GitHub action](https://github.com/Zomzog/changelog-checker) makes sure all pull requests that meet the above criteria have a change to the changelog in the files changed.

> Note: For further details on good practice for CHANGELOG entries see the [`keepachangelog` docs](https://keepachangelog.com/en/1.0.0/)

#### Update Issue When PR Linked

This job will only run if the [Verify Linked Issue](# Verify-Linked-Issue-Job) job has run successfully. It has a number of steps that use GitHub GraphQL queries to retrieve and update data:

1. **Get the project data**: gets variables such as the current sprint and specific Kanban board column names.
1. **Get linked issue `nodeid`**: gets the `nodeID` of the **most recent** issue that has been linked to the pull request.
1. **Add issue to project**: adds the linked issue to the project board, if not already present.
1. **Set issue project status fields**: updates the issue fields; sets the status to `In Progress` and sprint to `@current`. If the issue is already in `Waiting Review` the issue will remain with this status.

#### Skip Changelog And Issue Checks

This job will run on all **non-draft** pull requests that are not created by the Renovate Bot. Adding and removing the `no-changelog` label from a non-draft pull request will have the following effect:

* **adding** the `no-changelog` label will **start running this job**
* **removing** the `no-changelog` label will **stop this job from running**

This job has a number of steps that use GitHub GraphQL queries to retrieve and update data:

1. **Get the project data**: gets variables such as the current sprint and specific Kanban board column names.
1. **Add `pr` to project**: adds the pull request to the project board so the team knows it exists.
1. **Set `pr` project status fields**: updates the issue fields; sets the status to `Waiting Review` and sprint to `@current`.
