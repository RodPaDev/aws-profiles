import chalk from "chalk";
import { Command } from "commander";
import inquirer from "inquirer";

import packageJson from "../package.json";

import AWSProfiler from "./AwsProfiler";
import { createSpaces } from "./util";
import { Print } from "./util/Print";

const defaultPath = AWSProfiler.getDefaultConfigPath();

const cli = new Command();

cli
  .name("awsii")
  .description(packageJson.description)
  .version(packageJson.version);

cli.option(
  "-p, --path <path>",
  "Custom path to the aws-cli credentials file.",
  defaultPath
);

const global = cli.opts();

cli
  .command("list")
  .description("List all profiles")
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path);
    awsProfiler.findCurrentProfile();
    const profiles = awsProfiler.createProfileList();
    Print.info("List of aws profiles:");
    let index = 0;
    for (let profile of profiles) {
      if (index === 0) {
        Print.info(createSpaces(2), chalk.cyan(`- ${profile} [Default]`));
      } else {
        Print.info(createSpaces(2), chalk.white(`- ${profile}`));
      }
      index += 1;
    }
  });

cli
  .command("select")
  .description("Select a profile to select as the default aws profile")
  .argument("<profile-name>", "Profile name")
  .action((profileName, options, command) => {
    const awsProfiler = new AWSProfiler(global.path);
    const profile = awsProfiler.findCurrentProfile();
    if (profile === profileName) {
      Print.info(
        chalk.white.bold("Default profile is already:"),
        chalk.cyan(profileName)
      );
    } else {
      switchProfile(awsProfiler, profileName);
    }
  });

cli
  .command("where")
  .description("Print the path to the aws credentials file")
  .action(() => {
    Print.info(chalk.greenBright("Path:"), chalk.white(defaultPath));
  });

cli
  .command("choose")
  .description("Choose a profile from an interactive shell")
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path);
    awsProfiler.findCurrentProfile();
    const [current, ...choices] = awsProfiler.createProfileList();
    Print.info("Current default aws profile: ", chalk.cyan(current));
    inquirer
      .prompt([
        {
          type: "list",
          name: "profileName",
          message: "Choose a new default profile:",
          choices,
        },
      ])
      .then(({ profileName }) => {
        Print.log(' ')
        switchProfile(awsProfiler, profileName);
      });
  });

function switchProfile(awsProfiler, newProfile) {
  awsProfiler.selectProfile(newProfile);
  awsProfiler.saveProfile();
  Print.info(
    chalk.white("Default profile switched to: "),
    chalk.cyan(newProfile)
  );
}

cli.parse();
