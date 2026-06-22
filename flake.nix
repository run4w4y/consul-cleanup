{
  description = "consul-cleanup development shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/cb9a96f23c491c081b38eab96d22fa958043c9fa";
    nixpkgs-act.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, nixpkgs-act, flake-utils }: flake-utils.lib.eachDefaultSystem (
    system:
    let
      goVersion = 22; # 1.xx go version

      goOverlay = final: prev: {
        go = final."go_1_${toString goVersion}";
      };

      pkgs = import nixpkgs {
        inherit system;
        config.allowUnfree = true;
        overlays = [ goOverlay ];
      };

      actPkgs = import nixpkgs-act {
        inherit system;
        config.allowUnfree = true;
      };
    in
    with pkgs;
    {
      devShells.default = mkShell {
        packages = [
          actPkgs.act
          actPkgs.opentofu
          go
          gotools
          golangci-lint
        ];
      };
    }
  );
}
