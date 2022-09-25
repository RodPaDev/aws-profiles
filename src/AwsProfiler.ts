import { homedir } from "os";
import { resolve } from "path";
import { existsSync, PathLike } from "fs";
import { createIniConfig, readConfig, writeConfig } from "./util/IniConfig";
import { Print } from "./util/Print";
import { exit } from "process";

type Config = {
  [key: string]: {
    aws_access_key_id: string;
    aws_secret_access_key: string;
  };
};

class AWSProfiler {
  configPath: PathLike = "";
  config: Config = {};
  defaultProfileName: string = "";

  constructor(awsCredentialsPath: PathLike) {
    this.configPath = awsCredentialsPath;
    if (!existsSync(this.configPath)) {
      Print.error("AWS credentials file could not be found.");
      exit(1);
    }

    this.config = readConfig(this.configPath);
  }

  findCurrentProfile(): string {
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
    return "";
  }

  createProfileList(appendToCurrent: string = ""): string[] {
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

  selectProfile(name: string): string {
    this.config.default = this.config[name];
    return `Default profile set to '${name}'`;
  }

  saveProfile(): void {
    const iniConfig = createIniConfig(this.config);
    writeConfig(this.configPath, iniConfig);
  }

  static getDefaultConfigPath(): string {
    return resolve(homedir(), ".aws", "credentials");
  }
}

export default AWSProfiler;
