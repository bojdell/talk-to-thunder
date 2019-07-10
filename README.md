# talk-to-thunder

This repo contains a sample app that is a companion to an in-person workshop. The primary goal of the workshop is to learn about new technologies in a fun, hands-on way (particularly interns/new grads). The project is built on the following technologies:

- GraphQL
- Go / Thunder
- React / Typescript
- Docker
- MySQL
- GPT-2 language model

Content from previous workshops can be found at [#talktothunder](https://twitter.com/hashtag/talktothunder) on Twitter. The workshop slides can also be found [here](TODO).

## Overview

The workshop is split into the following components:

1. Intro / GraphQL Basics (~15 mins)
2. Warmup (~5 mins)
3. Installation / Setup (~20 mins)
4. Building a Product with GraphQL (~60 mins)
   - Warmup (~5 mins)
   - Challenge 1 (~20 mins)
   - Challenge 2 (~20 mins)
   - Challenge 3 (~20 mins)

## 2) Warmup

Go to GitHub's public [GraphQL API Explorer](https://developer.github.com/v4/explorer/):

- Can you find out when the [thunder](https://github.com/samsarahq/thunder) repo was created?
- Can you fetch the list of languages used in it?
- Advanced / Optional: If you're feeling up to it, can you use the GraphQL API to star it?
  - Hint: you might need to first find the ID of the repository 👀

## 3) Installation / Setup (~20 mins)

- Clone this project in your terminal by running `git clone https://github.com/bojdell/talk-to-thunder.git`. This will copy all the code to your computer under the directory `talk-to-thunder`. If you don't have [Git](https://git-scm.com/downloads), you will need to install it.
- Run `./setup.sh` to install various dependencies for the project.
- Run each of the following in its own terminal window:
  - `docker-compose -f db/docker-compose.yml up`
    - Note: if you get port conflicts, make sure to shut down any other MySQL instances running on your machine.
    - Run `migrate -database 'mysql://root:@tcp(127.0.0.1:3307)/talktothunder' -path ./db/migrations up` in a new teminal window.
  - `yarn start`
  - `go run go/src/talktothunder/gqlserver/main.go`
  - `docker-compose -f db/docker-compose.yml up`
  - `docker pull` TODO

#### Other Tips

- Have a lightweight text editor available (we strongly recommend [Visual Studio Code](https://code.visualstudio.com/))
- Be semi-comfortable navigating your terminal of choice

## 4) Building a Product with GraphQL (~60 mins)

TODO

### Challenges 1-3

Follow the instrutions on each page to solve the challenge. Again remember that you can use the GraphiQL to test and debug things. One participant will share her/his solution to each challenge.

### Feedback

You can provide feedback on the workshop here: TODO. Please help us make it better for future sessions!
