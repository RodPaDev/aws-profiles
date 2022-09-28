import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'fs'
import os from 'os'
import path from 'path'

type CliConfig = {
  name: string
  version: string
  description: string
  reservedNames: string[]
  customPath: string
}

class CliConfigManager {
  path: string
  config: CliConfig
  constructor() {
    this.path = path.join(os.homedir(), '.aws-profiles', 'config.json')
    this.config = this.load()
  }

  setCustomPath(path: string) {
    this.config.customPath = path
  }

  save(): void {
    writeFileSync(this.path, JSON.stringify(this.config, null, 2))
  }

  load(): CliConfig {
    if (!existsSync(this.path)) {
      const pathObj = path.parse(this.path)
      mkdirSync(pathObj.dir, { recursive: true })
      writeFileSync(this.path, JSON.stringify(this.config ?? {}, null, 2), {
        encoding: 'utf-8'
      })
    }
    return <CliConfig>JSON.parse(readFileSync(this.path, { encoding: 'utf-8' }))
  }
}

export default CliConfigManager
