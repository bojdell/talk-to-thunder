# talk-to-thunder

This repo contains a sample app that is a companion to an in-person workshop. The primary goal of the workshop is to learn about new technologies in a fun, hands-on way (particularly interns/new grads). The project is built on the following technologies:

- GraphQL
- Go / [Thunder](https://github.com/samsarahq/thunder)
- React / Typescript
- Docker
- MySQL
- [GPT-2 language model](https://openai.com/blog/better-language-models/)

Content from previous workshops can be found at [#talktothunder](https://twitter.com/hashtag/talktothunder) on Twitter. The workshop slides can also be found [here](https://docs.google.com/presentation/d/1UEk_hkmv0Jgmxhq36E0C77KOSikTbsOFwWySNpj13RY/edit?usp=sharing). The app is inspired by https://talktotransformer.com

![image](https://user-images.githubusercontent.com/3486165/61419524-74945000-a8b3-11e9-9e84-58821622576e.png)

## Overview

The workshop is split into the following components:

1. Intro / GraphQL Basics (~15 mins)
2. Installation / Setup (~20 mins)
3. Building a Product with GraphQL (~45 mins)
   - Warmup (~5 mins)
   - Challenge 1 (~10 mins)
   - Challenge 2 (~10 mins)
   - Challenge 3 (~20 mins)

## 1) Intro / GraphQL Basics (~15 mins)

[Slides](https://docs.google.com/presentation/d/1UEk_hkmv0Jgmxhq36E0C77KOSikTbsOFwWySNpj13RY/edit?usp=sharing)

## 2) Installation / Setup (~20 mins)

- Clone this project in your terminal by running `git clone https://github.com/bojdell/talk-to-thunder.git`. This will copy all the code to your computer under the directory `talk-to-thunder`. If you don't have [Git](https://git-scm.com/downloads), you will need to install it.
- `cd talk-to-thunder`
- Run each of the following in a new terminal tab:
  - `docker pull bojdell/talktothunder:gpt-2.1` (takes a while)
  - `./setup.sh` to install various dependencies for the project.
- Run each of the following in a new terminal tab:
  - `docker-compose -f db/docker-compose.yml up`
    - Note: if you get port conflicts, make sure to shut down any other MySQL instances running on your machine.
  - `migrate -database 'mysql://root:@tcp(127.0.0.1:3307)/talktothunder' -path ./db/migrations up` (you can reuse this window after it runs)
    - If you ever need to login to the database, you can run `mysql -h 127.0.0.1 --port=3307 -uroot talktothunder`
  - `yarn start`
  - `go run go/src/talktothunder/gqlserver/main.go`
  - `docker run --rm -it -v /tmp/talktothunder:/tmp/talktothunder:rw bojdell/talktothunder:gpt-2.1 python3 src/incremental.py --top_k 60 --model_name 345M`

#### Other Tips

- Have a lightweight text editor available (we strongly recommend [Visual Studio Code](https://code.visualstudio.com/))
- Be semi-comfortable navigating your terminal of choice

## 3) Building a Product with GraphQL (~45 mins)

### Warmup

Go to GitHub's public [GraphQL API Explorer](https://developer.github.com/v4/explorer/):

- Can you find out when the [thunder](https://github.com/samsarahq/thunder) repo was created?
- Can you fetch the list of languages used in it?
- Advanced / Optional: If you're feeling up to it, can you use the GraphQL API to star it?
  - Hint: you might need to first find the ID of the repository 👀

### Challenges 1-3

Follow the instrutions on each page to solve the challenge. Again remember that you can use the GraphiQL to test and debug things. One participant will share her/his solution to each challenge.

- Challenge 1: client/components/all_snippets.tsx
- Challenge 2: client/components/snippet.tsx
- Challenge 3a: go/src/talktothunder/gqlserver/main.go
- Challenge 3b: client/components/share_buttons.tsx
- Bonus: ???

### Feedback

You can provide feedback on the workshop [here](https://forms.gle/yKdVi6gv7Vt4QSLj9). Please help us make it better for future sessions!
