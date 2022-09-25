import { homedir } from "os";
import { resolve } from "path";
import { existsSync } from "fs";
import { createIniConfig, readConfig, writeConfig } from "./util/IniConfig";
import { Print } from "./util/Print";
import { exit } from "process";

class AWSProfiler {
  configPath = null;
  config = null;
  defaultProfileName = null;

  constructor(awsCredentialsPath = null) {
    this.configPath = awsCredentialsPath;
    if (!existsSync(this.configPath)) {
      Print.error("AWS credentials file could not be found.");
      exit(1);
    }

    this.config = readConfig(this.configPath);
  }

  findCurrentProfile() {
    for (let profile in this.config) {
      if (
        profile !== "default" &&
        this.config[profile].aws_access_key_id ===
          this.config.default.aws_access_key_id
      ) {
        this.defaultProfileName = profile;
        return profile;
      }
    }
    return null;
  }

  createProfileList(appendToCurrent = "") {
    const profiles = [];
    for (let profile of Object.keys(this.config)) {
      if (profile !== "default") {
        if (profile === this.defaultProfileName) {
          profiles.unshift(profile + appendToCurrent);
        } else {
          profiles.push(profile);
        }
      }
    }
    return profiles;
  }

  selectProfile(name) {
    this.config.default = this.config[name];
    return `Default profile set to '${name}'`;
  }

  saveProfile() {
    const iniConfig = createIniConfig(this.config);
    writeConfig(this.configPath, iniConfig);
  }

  static getDefaultConfigPath() {
    return resolve(homedir(), ".aws", "credentials");
  }
}

export default AWSProfiler;
