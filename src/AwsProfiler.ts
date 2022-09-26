import { homedir } from 'os'
import { resolve } from 'path'
import { existsSync, PathLike } from 'fs'
import {
  Config,
  createIniConfig,
  readConfig,
  writeConfig
} from './utils/IniConfig'
import { Print } from './utils/Print'
import { exit } from 'process'

class AWSProfiler {
  configPath: PathLike
  config: Config
  defaultProfileName: string

  constructor(awsCredentialsPath: PathLike) {
    this.configPath = awsCredentialsPath
    if (!existsSync(this.configPath)) {
      Print.error('AWS credentials file could not be found.')
      exit(1)
    }

    this.config = readConfig(this.configPath)
    this.defaultProfileName = ''
  }

  findCurrentProfile(): string {
    for (const profile in this.config) {
      if (
        profile !== 'default' &&
        this.config[profile].aws_access_key_id ===
          this.config.default.aws_access_key_id
      ) {
        this.defaultProfileName = profile
        return profile
      }
    }
    return ''
  }

  createProfileList(appendToCurrent = ''): string[] {
    const profiles = []
    for (const profile of Object.keys(this.config)) {
      if (profile !== 'default') {
        if (profile === this.defaultProfileName) {
          profiles.unshift(profile + appendToCurrent)
        } else {
          profiles.push(profile)
        }
      }
    }
    return profiles
  }

  selectProfile(name: string): string {
    this.config.default = this.config[name]
    return `Default profile set to '${name}'`
  }

  saveProfile(): void {
    const iniConfig = createIniConfig(this.config)
    writeConfig(this.configPath, iniConfig)
  }

  static getDefaultConfigPath(): string {
    return resolve(homedir(), '.aws', 'credentials')
  }
}

export default AWSProfiler
