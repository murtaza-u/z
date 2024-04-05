{
  description = "Murtaza Udaipurwala's Go Monolith";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.05";
    flake-utils.url = "github:numtide/flake-utils/v1.0.0";
  };
  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        formatter = pkgs.nixpkgs-fmt;
        packages = {
          default = pkgs.buildGoModule {
            pname = "z";
            version = "0.2.0";
            src = ./.;
            vendorSha256 = "sha256-U5c5jagaRW+lq0jsxWnngdcU2YYEl8Jn2Phuuq5vdzs=";
            CGO_ENABLED = 0;
            subPackages = [ "cmd/z" ];
            nativeBuildInputs = [ pkgs.installShellFiles ];
            postInstall = ''
              for shell in bash zsh; do
                installShellCompletion --$shell ./completion/$shell/z
              done
            '';
          };
        };
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            go
            go-tools
            gopls
          ];
        };
      });
}
