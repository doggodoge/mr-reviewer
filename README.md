# MR Reviewer

I was growing a bit sick of plodding through GitLab's UI to find MRs that are
ready for review, so made this quick dirty TUI to access non-draft (ready to review)
MRs quickly.

<img width="1037" alt="image" src="">

## Features

- Attractive TUI Interface
- Quickly drill down and find a MR from the comfort of your command line
- Toggle draft PRs by pressing <kbd>d</kbd>.
- Slap enter on a MR in the tool, and it opens in your browser!

## Build

```shell
go mod tidy
go build
```

A statically linked binary called `mr-reviewer` will appear in the repo. You
can copy it to `/usr/local/bin` to instantly use it!

## Config

You will need a tiny bit of configuration to add your bearer token.

- CD into `~/.config`
- Create a folder called `mr-reviewer`
- Create a `config.json` inside this new folder
- Open the file and add the boilerplate below
- Fill in the gitlab base path, your gitlab token, and the list of repositories
  you want to track

```json
{
  "gitlab_base_path": "<Your GitLab instances base url, eg. 'https://gitlab.com'>",
  "gitlab_token": "<Your gitlab access token>",
  "repositories": [
    {
      "name": "<repo name>",
      "description": "<repo description>",
      "route": "<route/to/project>"
    }
  ]
}
```

### How to Generate an Access Token

- Open GitLab
- Click your profile icon
- Navigate to "Preferences"
- Click "Access Tokens" on the left sidebar
- Generate a new token with any name, and the scope `read_api`
- Your token will appear at the top of the page, below the search bar
- Copy it into the json file described abov

## Needs work

- The code is pretty messy ;-;
- Absolutely no tests
- Missing help messages at the bottom for `d` and `backspace`. Shouldn't be too
  challenging to add.
- A UI for creating the configuration for repositories would be nice.
