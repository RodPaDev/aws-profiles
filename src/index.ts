import chalk from 'chalk'
import { Command } from 'commander'
import inquirer from 'inquirer'
import { exit } from 'process'

import AWSProfiler from './AwsProfiler'
import { createSpaces } from './utils'
import { Print } from './utils/Print'

// const defaultDirPath = AWSProfiler.getDefaultConfigDirPath()
const defaultDirPath =
  '/Users/rodpadev/Desktop/projects/aws-profile-switcher/tests/config'
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
    const originalProfileName =
      awsProfiler.getCaseInsensitiveProfileName(profileName)

    if (!originalProfileName) {
      Print.error('Profile does not exist')
      exit(1)
    }

    if (profile === originalProfileName) {
      Print.info(
        chalk.white.bold('Default profile is already:'),
        chalk.cyan(originalProfileName)
      )
    } else {
      switchProfile(awsProfiler, originalProfileName)
    }
  })

cli
  .command('current')
  .description("Display the current Default Profile's name")
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path)
    awsProfiler.findCurrentProfile()
    Print.info(
      'Current default aws profile: ',
      chalk.cyan(awsProfiler.defaultProfileName)
    )
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

cli
  .command('add')
  .description('Add a new aws profile')
  .action(() => {
    const awsProfiler = new AWSProfiler(global.path)
    awsProfiler.findCurrentProfile()
    inquirer
      .prompt([
        {
          type: 'input',
          name: 'name',
          message: 'Enter a profile name: ',
          validate: name =>
            awsProfiler.getCaseInsensitiveProfileName(name)
              ? 'Profile name is already in use'
              : true
        },
        {
          type: 'input',
          name: 'aws_access_key_id',
          message: 'aws_access_key_id: '
        },
        {
          type: 'password',
          name: 'aws_secret_access_key',
          message: 'aws_secret_access_key: '
        },
        {
          type: 'input',
          name: 'region',
          message: 'region: '
        },
        {
          type: 'input',
          name: 'output',
          message: 'output: '
        },
        {
          type: 'confirm',
          name: 'setDefault',
          message: 'Set this profile as default? ',
          default: false
        }
      ])
      .then(
        ({
          name,
          aws_access_key_id,
          aws_secret_access_key,
          region,
          output,
          setDefault
        }) => {
          const credentials = { aws_access_key_id, aws_secret_access_key }
          const config = { region, output }
          awsProfiler.addProfile(credentials, config, name)
          awsProfiler.saveProfile()

          if (setDefault) {
            switchProfile(awsProfiler, name, false)
            Print.info(
              'New aws profile added and set as default:',
              chalk.cyan(chalk.cyan(name))
            )
          } else {
            Print.info('New aws profile added:', chalk.cyan(chalk.cyan(name)))
          }
        }
      )
  })

function switchProfile(
  awsProfiler: AWSProfiler,
  profileName: string,
  print = true
) {
  awsProfiler.setDefaultProfile(profileName)
  awsProfiler.saveProfile()
  if (print) {
    Print.info(
      chalk.white('Default profile switched to: '),
      chalk.cyan(profileName)
    )
  }
}

cli.parse()
