#!/bin/bash
# setup.sh configures a local development environment.

cat <<EOF
Welcome to setup.sh! This script is going to setup your local development
environment. If at any point the script breaks, the first step to debugging it
is to run the script again. If you have a better fix, please submit a PR.

EOF

# # Set the working directory to the root of the repository.
# cd ${0%/*}/..ls .

# Figure out what OS we're running.
case $(uname -s) in
Darwin) OS=darwin ;;
Linux) OS=linux ;;
*) echo "Unsupport OS $(uname -s)"; exit 1 ;;
esac

# Helper function to echo a command before running it.
run() {
  echo "$@" && "$@"
}

# Fail on errors or unassigned variables, print all lines of the script.
set -eu

# echo Installing OS-wide binaries.
# case $OS in
# darwin)
#   echo Checking if brew installed.
#   if ! which brew >/dev/null 2>&1; then
#     echo Installing brew.
#     run ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
#   fi

#   echo Updating brew formulas.
#   run brew update
#   echo Ensuring there is no manually setup Docker around.
#   if [ -d /Applications/Docker.app ] && [ ! -d /usr/local/Caskroom/docker ]; then
# cat <<EOF
# setup.sh manages Docker on Mac using Homebrew. If you have previously manually
# installed Docker, you need to remove Docker first. Please stop Docker (see the little
# whale icon in the menu bar) and remove /Applications/Docker.app to continue. Don't
# worry, your containers will not be deleted in this process.

# This script will not continue while /Applications/Docker.app exists but
# /usr/local/Caskroom/docker does not.
# EOF
#     exit 1
#   fi

#   echo Installing and updating brew binaries.

#   run brew bundle -v --file=- <<-EOF
#     brew "jq"
#     brew "direnv"
#     brew "git"
#     cask "docker"
# EOF

# if ! pgrep "Docker" >/dev/null 2>&1 ; then
#   tput setaf 1 # Switch to red to get attention
#   printf "Please ensure Docker is running (see the little whale icon in the menu bar) by launching it from Launchpad..."
#   while ! pgrep "Docker" >/dev/null 2>&1 ; do
#     printf "."
#     sleep 2
#   done
#   tput setaf 2 # Thank the user in green :)
#   echo "thanks! Continuing"
#   tput sgr0 # Reset text formatting
# fi

# ;;
# linux)
#   echo Install binaries with apt-get.
#   run sudo apt-get update
#   run sudo apt-get install -y jq ffmpeg direnv inotify-tools python3.6 clang-format unzip python-pip apt-transport-https ca-certificates curl software-properties-common shellcheck

#   # https://docs.docker.com/install/linux/docker-ce/ubuntu/
#   curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
#   run sudo apt-key fingerprint 0EBFCD88
#   run sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

#   run sudo apt-get update
#   run sudo apt-get install -y docker-ce
#   run sudo python -m pip install -U docker-compose

#   run sudo apt-get install -y pkg-config libusb-1.0-0-dev libu2f-host-dev

#   # Add current user to docker group to allow docker without sudo.
#   sudo usermod -aG docker $USER

#   echo Checking if awscli is installed.
#   if ! which aws >/dev/null 2>&1; then
#     run sudo python -m pip install -U awscli
#   fi

# esac
# echo

echo Configuring shell for direnv
# If direnv wasn't setup, we'll prompt the user at the end to login and logout.
if ! env | grep -q DIRENV -; then
  direnv_note="1"
else
  direnv_note=""
fi

if [ ! -f ~/.bash_profile ]; then
  echo "No ~/.bash_profile found, creating ~/.bash_profile"
  touch ~/.bash_profile
fi

load_bashrc="if [ -f $HOME/.bashrc ]; then source $HOME/.bashrc; fi"

if ! grep -q bashrc ~/.bash_profile 2> /dev/null; then
  echo "Appending $load_bashrc to ~/.bash_profile"
  echo "$load_bashrc" >> ~/.bash_profile
fi

add_direnv() {
  if [ -f $1 ] && ! grep -q direnv $1 2> /dev/null; then
    echo Appending $2 to $1
    echo "$2" >> $1
  fi
}
touch ~/.bashrc
add_direnv ~/.bashrc 'eval "$(direnv hook bash)"'
add_direnv ~/.zshrc 'eval "$(direnv hook zsh)"'
add_direnv ~/.config/fish/config.fish 'direnv hook fish | source'
direnv allow
eval "$(direnv export bash)"
echo

echo Checking if we have go $GO_VERSION.
if [ ! -f opt/go$GO_VERSION/bin/go ]; then
  echo Downloading go.
  GO_DOWNLOAD_URL=https://golang.org/dl/go$GO_VERSION.$OS-amd64.tar.gz
  GO_TMP_PATH=/tmp/go$GO_VERSION.tar.gz
  run curl -fsSL "$GO_DOWNLOAD_URL" -o $GO_TMP_PATH
  run mkdir -p opt/go$GO_VERSION
  run tar -C opt/go$GO_VERSION -xzf $GO_TMP_PATH --strip-components 1
  run rm $GO_TMP_PATH

  echo Resetting go binaries.
  run rm -rf go/bin/*
fi
echo

echo Checking if we have node $NODE_VERSION.
if [ ! -f opt/node$NODE_VERSION/bin/node ]; then
  NODE_DOWNLOAD_URL=https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-$OS-x64.tar.gz
  NODE_TMP_PATH=/tmp/node$NODE_VERSION.tar.gz
  run curl -fsSL "$NODE_DOWNLOAD_URL" -o $NODE_TMP_PATH
  run mkdir -p opt/node$NODE_VERSION
  run tar -C opt/node$NODE_VERSION -xzf $NODE_TMP_PATH --strip-components 1
  run rm $NODE_TMP_PATH

  echo Resetting node_modules.
  run rm -rf node_modules/
fi
echo

echo Making sure we have yarn $YARN_VERSION installed.
run npm install -g yarn@$YARN_VERSION
echo

echo Running yarn to install node_modules.
run yarn
echo
