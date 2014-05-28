# PETAR

R="\[\033[m\]"
C1="\[\033[7;30m\]"
C2="\[\033[0;32m\]"
C3="\[\033[7;32m\]"
C4="\[\033[7;30m\]"

export PS1="▒◀︎${C1} bg:\j ${R}◀︎${C1} \u@\h ${R}◀︎${C4} \w ${R}${C4}${R}\n▒ ${C4}${R}"
tabs -4

# shopt -s autocd

# TMUX: set tmux window title to current directory
function tw {
	gp=`gps .`
	if [ -z "$gp" ]
	then
		printf '\033k…/%s\033\\' `basename \`pwd\``
	else
		printf '\033k%s/…/%s\033\\' `basename $gp` `basename \`pwd\``
	fi
}

# aliases
alias scpresume='rsync --partial --progress --rsh=ssh'
alias a='acme -a -F ~/aux/plan9/font/lucsans/typeunicode.9.font -l ~/acme.dump'
alias e='vim'
alias d='ls -FA'
alias v='d -lA'
alias rm='rm' # PLAN9
alias kk='kill -KILL'
alias l='less -S'
alias g='ack'
alias h='head'
alias t='tail'
alias c='cat'
alias gs=""  # disable gs command

alias c0m="cd ~/0/src/github.com/petar/maymounkov.io"
alias c0c="cd ~/0/src/github.com/gocircuit/circuit"

# HG aliases
alias hll='hg log | less'

# GIT aliases
alias gu='git add -u'
alias ga='git add'
alias ggs='git status'
alias ggm=' git commit -m '
alias gqam='git pull am master'
alias gpam='git push -u am master'
alias gqom='git pull origin master'
alias gpom='git push -u origin master'

# GO TOOLS
alias gd_='godoc -http=:6060 -index'
alias gd='godoc'
alias gb='go build'
alias gi='go install'
alias gt='go test'
alias ucgo='unset CGO_LDFLAGS; unset CGO_CFLAGS'

# JULIA
alias julia='/Applications/Julia-0.2.0.app/Contents/Resources/julia/bin/julia'

# go appengine
alias hd='~/appengine_sdk/godoc'

# CD aliases
alias cg='cd ~/go/src'
alias c0='cd ~/0/src'
alias ch='pushd ./ ; cd ~/go/src && hg sync ; popd'

# VAGRANT
alias vg=vagrant
alias dckr=docker

# PATH
export PATH=/usr/local/share/python:${HOME}/aux/protobuf/bin:${HOME}/aux/bin:/usr/local/bin:/usr/bin:/bin:/sbin:/usr/sbin:/usr/X11R6/bin:/usr/local/sbin
export PATH=/usr/local/git/bin:${PATH}

# GO
export GOOS=linux
export GOARCH=amd64
export GOMAXPROCS=10
export GOPATH=$HOME/0
export GOROOT=$HOME/go
# export GODEBUG="gctrace=2, schedtrace=X, scheddetail=1"
export PATH=$PATH:$HOME/0/bin:$GOROOT/bin

# CIRCUIT
declare -x CIRCUIT=/Users/petar/0/src/github.com/gocircuit/circuit/cmd/circuit/.circuit
declare -x CIRCUIT_HMAC=/Users/petar/0/src/github.com/gocircuit/circuit/cmd/circuit/.hmac

export VAGRANT_CWD=$GOPATH/src/github.com/gocircuit/circuit/misc/vagrant
# export VAGRANT_GO_ORIGIN=$HOME/vagrant/go
# export VAGRANT_GOCIRCUIT_ORIGIN=$HOME/vagrant/gocircuit

# PLAN9
PLAN9=$HOME/aux/plan9port export PLAN9
PATH=$PATH:$PLAN9/bin export PATH

# SHELL
export PAGER='less -S'
umask 022

# MARKING
export MARKPATH=$HOME/.marks
function jmp { 
    cd -P $MARKPATH/$1 2>/dev/null || echo "No such mark: $1"
}
function mrk { 
    mkdir -p $MARKPATH; ln -s $(pwd) $MARKPATH/$1
}
function umrk { 
    rm -i $MARKPATH/$1 
}
function mrks {
    ls -l $MARKPATH | sed 's/  / /g' | cut -d' ' -f9- | sed 's/ -/ -/g' && echo
}
