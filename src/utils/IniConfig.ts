import { encode, parse } from 'ini'
import { PathLike, readFileSync, writeFileSync } from 'fs'

export type Config = {
  [key: string]: {
    aws_access_key_id: string
    aws_secret_access_key: string
  }
}

export function readConfig(path: PathLike): Config {
  return parse(readFileSync(path, { encoding: 'utf-8' }))
}

export function writeConfig(path: PathLike, iniConfig: string): void {
  writeFileSync(path, iniConfig, { encoding: 'utf-8' })
}

export function createIniConfig(configObject: Config): string {
  return encode(configObject)
}
