with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "nix-cage-shell";
  buildInputs = [
    influxdb
    nsq
    go
    gocode
    godef
    dep
  ];
  shellHook = ''
    export GOPATH=~/projects
  '';
}
