{
  lib,
  buildGo126Module,
}:

buildGo126Module {
  pname = "csv2ynab";
  version = "0.0.1";

  src = lib.cleanSource ../.;

  vendorHash = "sha256-UhogdYYb+N2q2hbhu7XSbaqe7yZOTW4haL31n+Nt5K0=";

  meta = {
    description = "Import Belgian bank exports into YNAB";
    homepage = "https://github.com/victorjacobs/csv2ynab";
    license = lib.licenses.mit;
    mainProgram = "csv2ynab";
  };
}
