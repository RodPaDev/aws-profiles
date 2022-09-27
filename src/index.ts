import chalk from 'chalk'
import { Command } from 'commander'
import inquirer from 'inquirer'
import { exit } from 'process'

import AwsProfileManager from './AwsProfileManager'
import { createSpaces } from './utils'
import { Print } from './utils/Print'

import config from './config.json'

const defaultDirPath = AwsProfileManager.getDefaultConfigDirPath()
const cli = new Command()

cli.name(config.name).description(config.description).version(config.version)

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
    const awsProfiler = new AwsProfileManager(global.path)
    awsProfiler.findCurrentProfile()
    const profiles = awsProfiler.createProfileList()
    Print.info('List of aws profiles:')
    let index = 0
    for (const profile of profiles) {
      if (index === 0) {
        Print.log(createSpaces(2), chalk.cyan(`- ${profile} [Default]`))
      } else {
        Print.log(createSpaces(2), chalk.white(`- ${profile}`))
      }
      index += 1
    }
  })

cli
  .command('select')
  .description('Select a profile to select as the default aws profile')
  .argument('<profile-name>', 'Profile name')
  .action(profileName => {
    const awsProfiler = new AwsProfileManager(global.path)
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
    const awsProfiler = new AwsProfileManager(global.path)
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
    const awsProfiler = new AwsProfileManager(global.path)
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
        switchProfile(awsProfiler, profileName)
      })
  })

cli
  .command('add')
  .description('Add a new aws profile')
  .action(() => {
    const awsProfiler = new AwsProfileManager(global.path)
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
              chalk.cyan(name)
            )
          } else {
            Print.info('New aws profile added:', chalk.cyan(name))
          }
        }
      )
  })

cli
  .command('delete')
  .description('Delete a aws profile')
  .argument('<profile-name>', 'Profile name')
  .action(profileName => {
    const awsProfiler = new AwsProfileManager(global.path)
    if (!awsProfiler.getCaseInsensitiveProfileName(profileName)) {
      Print.error('Profile does not exist')
      exit(1)
    }
    inquirer
      .prompt([
        {
          type: 'confirm',
          name: 'confirm',
          message: chalk.yellow(
            'This action is irreversible. Are you sure you want to delete?'
          ),
          default: false,
          prefix: chalk.yellow('!')
        }
      ])
      .then(({ confirm }) => {
        awsProfiler.deleteProfile(profileName)
        if (confirm) {
          awsProfiler.saveProfile()
          Print.info('AWS Profile deleted:', chalk.cyan(profileName))
        } else {
          Print.info('AWS Profile was not deleted:', chalk.cyan(profileName))
        }
      })
  })

function switchProfile(
  awsProfiler: AwsProfileManager,
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
