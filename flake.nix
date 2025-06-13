{
  description = "Geminic";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        pname = "geminic";
        version = "0.3.2";
      in {
        packages = {
          default = pkgs.buildGoModule {
            inherit pname version;
            src = ./.;
            hash = "sha256-77h2H/mhw9YkAMasVMxoadJbmjmxUPLf/k6tM+cOZcs=";
            vendorHash = "sha256-77h2H/mhw9YkAMasVMxoadJbmjmxUPLf/k6tM+cOZcs=";
          };
        };

        devShells = {
          default = pkgs.mkShell {
            buildInputs = [ pkgs.go pkgs.gopls pkgs.delve pkgs.gotools ];
          };
        };

        apps = {
          default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/${pname}";
          };
        };
      });
}
