import { encode, parse } from 'ini'
import { readFileSync, writeFileSync } from 'fs'

type GenericMap = { [key: string]: unknown }

export function readConfig(path: string): GenericMap {
  return parse(readFileSync(path, { encoding: 'utf-8' }))
}

export function writeConfig(path: string, iniConfig: string): void {
  writeFileSync(path, iniConfig, { encoding: 'utf-8' })
}

export function createIniConfig(configObject: GenericMap): string {
  return encode(configObject)
}
