with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "gluttony-shell";
  buildInputs = [
    influxdb
    go
    gocode
    glide
    godef
  ];
  shellHook = ''
    export GOPATH=~/projects
  '';
}
