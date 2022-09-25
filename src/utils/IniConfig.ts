import { encode, parse } from "ini";
import { PathLike, readFileSync, writeFileSync } from "fs";

export function readConfig(path: PathLike): any {
  return parse(readFileSync(path, { encoding: "utf-8" }));
}

export function writeConfig(path: PathLike, iniConfig: string): void {
  writeFileSync(path, iniConfig, { encoding: "utf-8" });
}

export function createIniConfig(configObject: any): string {
  return encode(configObject);
}
