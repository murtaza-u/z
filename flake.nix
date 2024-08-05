{
  description = "Monolith Go Commander";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        packages = {
          default = pkgs.buildGoModule {
            pname = "z";
            version = "0.0.1";
            src = ./.;
            vendorHash = null;
            CGO_ENABLED = 0;
            subPackages = [ "cmd/z" ];
          };
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            nixpkgs-fmt
            nixd
            go
            go-tools
            gopls
          ];
        };
      });
}
