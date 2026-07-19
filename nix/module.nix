{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.csv2ynab;
  json = pkgs.formats.json { };
  settingsFile = json.generate "csv2ynab-settings.json" cfg.settings;
  runtimeConfig = "/run/csv2ynab/config.json";
in
{
  options.services.csv2ynab = {
    enable = lib.mkEnableOption "csv2ynab WebDAV server";

    package = lib.mkOption {
      type = lib.types.package;
      default = pkgs.callPackage ./package.nix { };
      defaultText = lib.literalExpression "pkgs.callPackage ./nix/package.nix { }";
      description = "The csv2ynab package to use.";
    };

    bindAddress = lib.mkOption {
      type = lib.types.str;
      default = "127.0.0.1:8080";
      description = "Address on which the WebDAV server listens.";
    };

    settings = lib.mkOption {
      inherit (json) type;
      default = { };
      example = {
        watch.patterns = [
          {
            pattern = "*.csv";
            budget_id = "00000000-0000-0000-0000-000000000000";
            account_id = "00000000-0000-0000-0000-000000000000";
          }
        ];
      };
      description = ''
        csv2ynab configuration written as JSON. Do not put the YNAB API key
        here; use {option}`ynabAPIKeyFile`.
      '';
    };

    ynabAPIKeyFile = lib.mkOption {
      type = lib.types.path;
      description = "File containing the YNAB API key.";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.csv2ynab = {
      description = "csv2ynab WebDAV server";
      wantedBy = [ "multi-user.target" ];
      wants = [ "network-online.target" ];
      after = [ "network-online.target" ];

      serviceConfig = {
        Type = "simple";
        User = "csv2ynab";
        Group = "csv2ynab";
        RuntimeDirectory = "csv2ynab";
        RuntimeDirectoryMode = "0700";
        LoadCredential = "ynab-api-key:${cfg.ynabAPIKeyFile}";
        ExecStartPre = pkgs.writeShellScript "csv2ynab-config" ''
          ${lib.getExe pkgs.jq} \
            --rawfile apiKey "$CREDENTIALS_DIRECTORY/ynab-api-key" \
            '.ynab.api_key = ($apiKey | rtrimstr("\n"))' \
            ${settingsFile} > ${runtimeConfig}
        '';
        ExecStart = "${lib.getExe cfg.package} serve --bind ${cfg.bindAddress} --config ${runtimeConfig}";
        Restart = "on-failure";
        RestartSec = 10;
        UMask = "0077";

        CapabilityBoundingSet = "";
        LockPersonality = true;
        NoNewPrivileges = true;
        PrivateDevices = true;
        PrivateTmp = true;
        ProtectClock = true;
        ProtectControlGroups = true;
        ProtectHome = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectSystem = "strict";
        RestrictAddressFamilies = [
          "AF_INET"
          "AF_INET6"
        ];
        RestrictNamespaces = true;
        RestrictRealtime = true;
        SystemCallArchitectures = "native";
      };
    };

    users.groups.csv2ynab = { };
    users.users.csv2ynab = {
      isSystemUser = true;
      group = "csv2ynab";
      description = "csv2ynab service user";
    };
  };
}
