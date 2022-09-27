import { homedir } from 'os'
import path from 'path'
import { existsSync } from 'fs'
import { createIniConfig, readConfig, writeConfig } from './utils/IniConfig'
import { Print } from './utils/Print'
import { exit } from 'process'

type CredentialsObject = {
  aws_access_key_id: string
  aws_secret_access_key: string
}

type ConfigObject = {
  region: string
  output: string
}

type CredentialsMap = {
  [key: string]: CredentialsObject
}

type ConfigMap = {
  [key: string]: ConfigObject
}

class AWSProfiler {
  configPath: {
    credenetials: string
    config: string
  }
  config: ConfigMap
  credentials: CredentialsMap
  defaultProfileName: string

  constructor(configDir: string) {
    const config = path.join(configDir, 'config')
    const credenetials = path.join(configDir, 'credentials')
    this.configPath = {
      config,
      credenetials
    }

    if (!existsSync(this.configPath.credenetials)) {
      Print.error('AWS credentials file could not be found.')
      exit(1)
    }
    if (!existsSync(this.configPath.config)) {
      Print.error('AWS config file could not be found.')
      exit(1)
    }

    this.config = <ConfigMap>readConfig(this.configPath.config)
    this.credentials = <CredentialsMap>readConfig(this.configPath.credenetials)
    this.defaultProfileName = ''
  }

  findCurrentProfile(): string {
    for (const profile in this.credentials) {
      if (
        profile !== 'default' &&
        this.credentials[profile].aws_access_key_id ===
          this.credentials.default.aws_access_key_id
      ) {
        this.defaultProfileName = profile
        return profile
      }
    }
    return ''
  }

  createProfileList(appendToCurrent = ''): string[] {
    const profiles = []
    for (const profile of Object.keys(this.credentials)) {
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

  setDefaultProfile(name: string): string {
    this.config.default = this.config[name]
    this.credentials.default = this.credentials[name]
    return `Default profile set to '${name}'`
  }

  saveProfile(): void {
    const iniConfig = createIniConfig(this.config)
    const iniCredentials = createIniConfig(this.credentials)
    writeConfig(this.configPath.config, iniConfig)
    writeConfig(this.configPath.credenetials, iniCredentials)
  }

  addProfile(
    credentials: CredentialsObject,
    config: ConfigObject,
    profileName: string
  ): void {
    this.config[profileName] = config
    this.credentials[profileName] = credentials
  }

  deleteProfile(profileName: string): void {
    delete this.config[profileName]
    delete this.credentials[profileName]
  }

  static getDefaultConfigDirPath(): string {
    return path.resolve(homedir(), '.aws')
  }
}

export default AWSProfiler
