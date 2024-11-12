{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go              # Go programming language
    pkgs.gopls           # Go language server for editor support
    pkgs.postgresql      # PostgreSQL database
    pkgs.podman          # Podman for container management
    pkgs.git             # Git for version control, if needed
  ];

  # Set up Go environment variables
  shellHook = ''
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH
    echo "Go, Gorm, PostgreSQL, and Podman environment loaded!"
  '';
}
