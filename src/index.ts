import chalk from 'chalk'
import { Command } from 'commander'
import inquirer from 'inquirer'

import AWSProfiler from './AwsProfiler'
import { createSpaces } from './utils'
import { Print } from './utils/Print'

const defaultDirPath = AWSProfiler.getDefaultConfigDirPath()
const cli = new Command()

cli
  .name('awsii')
  .description('CLI to set named aws profiles to the default aws profile')
  .version('0.0.1')

cli.option(
  '-p, --path <path>',
  'Custom path to the aws-cli config folder.',
  defaultDirPath
)

const global = cli.opts()

cli
  .command('list')
  .description('List all profiles')
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path)
    awsProfiler.findCurrentProfile()
    const profiles = awsProfiler.createProfileList()
    Print.info('List of aws profiles:')
    let index = 0
    for (const profile of profiles) {
      if (index === 0) {
        Print.info(createSpaces(2), chalk.cyan(`- ${profile} [Default]`))
      } else {
        Print.info(createSpaces(2), chalk.white(`- ${profile}`))
      }
      index += 1
    }
  })

cli
  .command('select')
  .description('Select a profile to select as the default aws profile')
  .argument('<profile-name>', 'Profile name')
  .action(profileName => {
    const awsProfiler = new AWSProfiler(global.path)
    const profile = awsProfiler.findCurrentProfile()
    if (profile === profileName) {
      Print.info(
        chalk.white.bold('Default profile is already:'),
        chalk.cyan(profileName)
      )
    } else {
      switchProfile(awsProfiler, profileName)
    }
  })

cli
  .command('where')
  .description('Print the path to the aws config folder.')
  .action(() => {
    Print.info(chalk.greenBright('Path:'), chalk.white(defaultDirPath))
  })

cli
  .command('choose')
  .description('Choose a profile from an interactive shell')
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path)
    awsProfiler.findCurrentProfile()
    const [current, ...choices] = awsProfiler.createProfileList()
    Print.info('Current default aws profile: ', chalk.cyan(current))
    inquirer
      .prompt([
        {
          type: 'list',
          name: 'profileName',
          message: 'Choose a new default profile:',
          choices
        }
      ])
      .then(({ profileName }) => {
        Print.log(' ')
        switchProfile(awsProfiler, profileName)
      })
  })

function switchProfile(awsProfiler: AWSProfiler, profileName: string) {
  awsProfiler.setDefaultProfile(profileName)
  awsProfiler.saveProfile()
  Print.info(
    chalk.white('Default profile switched to: '),
    chalk.cyan(profileName)
  )
}

cli.parse()
