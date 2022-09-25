import { encode, parse } from "ini";
import { readFileSync, writeFileSync } from "fs";

export function readConfig(path) {
  return parse(readFileSync(path, { encoding: "utf-8" }));
}

export function writeConfig(path, iniConfig) {
  writeFileSync(path, iniConfig, { encoding: "utf-8" });
}

export function createIniConfig(configObject) {
  return encode(configObject);
}
