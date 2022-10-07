import chalk from 'chalk'
import { Command } from 'commander'
import inquirer from 'inquirer'
import { exit } from 'process'

import AwsProfileManager from './AwsProfileManager'
import { createSpaces } from './utils'
import { Print } from './utils/Print'

import config from './config.json'
import CliConfigManager from './utils/CliConfig'
import { existsSync } from 'fs'

const defaultDirPath = AwsProfileManager.getDefaultConfigDirPath()

const cli = new Command()
const cliConfigManager = new CliConfigManager()

cli.name(config.name).description(config.description).version(config.version)

cli.option(
  '-p, --path <path>',
  'Custom path to the aws-cli config folder.',
  defaultDirPath
)
cli.parse()

const global = cli.opts()

if (global.path !== defaultDirPath && existsSync(global.path)) {
  cliConfigManager.setCustomPath(global.path)
  cliConfigManager.save()
}

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

    if (originalProfileName === 'default') {
      Print.error('Cannot select', chalk.cyan('default'), 'profile')
      exit(1)
    }

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
          validate: name => {
            const [isValid, validFormat] = isNameValid(name)
            const originalProfileName =
              awsProfiler.getCaseInsensitiveProfileName(name)
            if (!isValid) {
              return 'Profile name must match: ' + chalk.magenta(validFormat)
            }
            if (originalProfileName) {
              return 'Profile already exists'
            }
            if (originalProfileName?.toLocaleLowerCase() === 'default') {
              return 'Cannot add profile named: ' + chalk.cyan('default')
            }
            return true
          }
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
          name: 'setAsDefault',
          message: 'Set this profile as default? '
        }
      ])
      .then(
        ({
          name,
          aws_access_key_id,
          aws_secret_access_key,
          region,
          output,
          setAsDefault
        }) => {
          const credentials = { aws_access_key_id, aws_secret_access_key }
          const config = { region, output }
          awsProfiler.upsertProfile(credentials, config, name)
          awsProfiler.saveProfile()

          if (setAsDefault) {
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
    const originalProfileName =
      awsProfiler.getCaseInsensitiveProfileName(profileName)
    if (!originalProfileName) {
      Print.error('Profile does not exist')
      exit(1)
    }
    if (originalProfileName === 'default') {
      Print.error('Cannot delete', chalk.cyan('default'), 'profile')
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

cli
  .command('edit')
  .description('Edit an existing aws profile')
  .action(() => {
    const awsProfiler = new AwsProfileManager(global.path)
    awsProfiler.findCurrentProfile()
    inquirer
      .prompt([
        {
          type: 'input',
          name: 'name',
          message: 'Enter a profile name: ',
          validate: name => {
            const [isValid, validFormat] = isNameValid(name)
            if (!isValid) {
              return 'Profile name must match: ' + chalk.magenta(validFormat)
            }
            const originalProfileName =
              awsProfiler.getCaseInsensitiveProfileName(name)
            if (!originalProfileName) {
              return 'Profile does not exist'
            }
            if (originalProfileName.toLocaleLowerCase() === 'default') {
              return 'Cannot edit default profile'
            }
            return true
          }
        },
        {
          type: 'confirm',
          name: 'isRename',
          message: 'Do you wish to rename the profile?'
        },
        {
          when: ({ isRename }) => isRename,
          type: 'input',
          name: 'newName',
          message: 'New profile name:',
          validate: name => {
            const originalProfileName =
              awsProfiler.getCaseInsensitiveProfileName(name)
            if (originalProfileName) {
              return 'Profile name is already in use'
            }
            if (originalProfileName?.toLocaleLowerCase() === 'default') {
              return 'Cannot rename to default'
            }
            return true
          }
        },
        {
          type: 'input',
          name: 'aws_access_key_id',
          message: 'aws_access_key_id: ',
          suffix: chalk.italic.grey('leave empty to skip edit'),
          default: (answers: { [key: string]: string }) => {
            return awsProfiler.credentials[answers.name].aws_access_key_id
          }
        },
        {
          type: 'password',
          name: 'aws_secret_access_key',
          message: 'aws_secret_access_key: ',
          suffix: chalk.italic.grey('leave empty to skip edit'),
          default: (answers: { [key: string]: string }) => {
            return awsProfiler.credentials[answers.name].aws_secret_access_key
          }
        },
        {
          type: 'input',
          name: 'region',
          message: 'region: ',
          suffix: chalk.italic.grey('leave empty to skip edit'),
          default: (answers: { [key: string]: string }) => {
            return awsProfiler.config[answers.name].region
          }
        },
        {
          type: 'input',
          name: 'output',
          message: 'output: ',
          suffix: chalk.italic.grey('leave empty to skip edit'),
          default: (answers: { [key: string]: string }) => {
            return awsProfiler.config[answers.name].output
          }
        },
        {
          type: 'confirm',
          name: 'isSave',
          message: 'Save the edits for this profile?'
        },
        {
          when: ({ isSave }) => isSave,
          type: 'confirm',
          name: 'setAsDefault',
          message: 'Set this profile as default? '
        }
      ])
      .then(
        ({
          name,
          isRename,
          newName,
          aws_access_key_id,
          aws_secret_access_key,
          region,
          output,
          setAsDefault,
          isSave
        }) => {
          const credentials = { aws_access_key_id, aws_secret_access_key }
          const config = { region, output }
          if (!isSave) {
            Print.warn('Edits not saved for aws profile:', chalk.cyan(name))
            exit(0)
          }
          if (isRename) {
            awsProfiler.deleteProfile(name)
            awsProfiler.upsertProfile(credentials, config, newName)
            Print.info(
              'Profile renamed:',
              chalk.yellow(name),
              '=>',
              chalk.cyan(newName)
            )
          } else {
            awsProfiler.upsertProfile(credentials, config, name)
          }
          awsProfiler.saveProfile()
          if (setAsDefault) {
            switchProfile(awsProfiler, name, false)
            Print.info(
              'Edits saved for aws profile and set as default:',
              chalk.cyan(isRename ? newName : name)
            )
          } else {
            Print.info(
              'Edits saved for aws profile:',
              chalk.cyan(isRename ? newName : name)
            )
          }
        }
      )
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

function isNameValid(name: string): [boolean, RegExp] {
  const validNameRegex = new RegExp(
    /^(?:@[a-z0-9-*~][a-z0-9-*._~]*\/)?[a-z0-9-~][a-z0-9-._~]*$/
  )
  return [validNameRegex.test(name), validNameRegex]
}

cli.parse()
