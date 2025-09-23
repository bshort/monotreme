# Monotreme

**Monotreme** is an open source, self-hosted platform designed to help you organize, manage, and share your most important links. Easily create customizable, human-readable shortcuts to streamline your link management. Use tags to categorize your links and share them easily with your team or publicly.

Monotreme is a fork of [Slash](https://github.com/monotrememarks/slash) and was created to allow extensions to Slash's core model to fit one user's bizarre requirements. Monotreme is probably not ready for widespread use and if you're looking for a self-hosted bookmarking application then just use [Slash](https://github.com/monotrememarks/slash). It's better!


## Background

In today's workplace, essential information is often scattered across the cloud in the form of links. We understand the frustration of endlessly searching through emails, messages, and websites just to find the right link. Links are notorious for being unwieldy, complex, and easily lost in the shuffle. Remembering and sharing them can be a challenge.

That's why we developed Monotreme (forked from [Slash](https://github.com/monotrememarks/slash)), a solution that transforms these links into easily accessible, discoverable, and shareable shortcuts(e.g., `s/shortcut`). Say goodbye to link chaos and welcome the organizational ease of Monotreme into your daily online workflow.

## Features

- Create customizable `s/` short links for any URL.
- Share short links public or only with your teammates.
- View analytics on link traffic and sources.
- Easy access to your shortcuts with browser extension.
- Share your shortcuts with Collection to anyone, on any browser.
- Open source self-hosted solution.

## Deploy with Docker in seconds

```bash
docker run -d --name monotreme -p 5231:5231 -v ~/.monotreme/:/var/opt/monotreme bshort/monotreme:latest
```

Learn more in [Self-hosting Monotreme with Docker](https://github.com/bshort/monotreme/blob/main/docs/install.md).

