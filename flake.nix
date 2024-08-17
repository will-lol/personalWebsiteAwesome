{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    templRepo.url = "github:a-h/templ/refs/tags/v0.2.543";
  };

  outputs = { self, nixpkgs, flake-utils, templRepo }:
    flake-utils.lib.eachDefaultSystem (system: 
      let 
        name = "personalWebsite";
	templ = templRepo.packages.${system}.templ;
	overlays = [];
	lib = nixpkgs.lib;
	pkgs = import nixpkgs { inherit system overlays; };
	goBuild = (pkgs.buildGoModule {
	  inherit name;
	  src = ./.;
	  CGO_ENABLED = 0;
	  ldflags = ["-s -w"];
	  vendorHash = "sha256-zAyiFzEHTmFwWSoqcfm4HMbc+CBkhr22Kar286JTipE=";
	  preBuild = ''
	    ${templ}/bin/templ generate
	  '';
	}).overrideAttrs (old: old // {GOOS = "linux"; GOARCH = "arm64"; });
      in
        {
	  packages = {
	    go = goBuild;
	  };
	  defaultPackage = goBuild;
	  devShell = pkgs.mkShell {
	    packages = [ pkgs.tailwindcss pkgs.entr pkgs.nodePackages.typescript-language-server pkgs.fd pkgs.nodejs_18 pkgs.awscli2 pkgs.aws-sam-cli pkgs.go pkgs.gopls pkgs.nodePackages.aws-cdk templ ];
	    shellHook = ''
	      export ENVIRONMENT=dev
	    '';
	  };
	}
    );
}
