const fs = require('fs')
const path = require('path')
const puppeteer = require('puppeteer-core')

const chromePath = 'C:/Program Files/Google/Chrome/Application/chrome.exe'
const userDataDir = path.join(__dirname, '.chrome-neon-profile')
const outputDir = path.join(__dirname, '.automation-output')
const outputFile = path.join(outputDir, 'neon-result.json')

function ensureDir(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true })
  }
}

async function clickFirst(page, selectors) {
  for (const selector of selectors) {
    try {
      if (selector.startsWith('xpath=')) {
        const xpath = selector.slice(6)
        await page.waitForSelector(`::-p-xpath(${xpath})`, { timeout: 2000 })
        await page.click(`::-p-xpath(${xpath})`)
        return true
      }
      await page.waitForSelector(selector, { timeout: 2000 })
      await page.click(selector)
      return true
    } catch {}
  }
  return false
}

async function typeFirst(page, selectors, value) {
  for (const selector of selectors) {
    try {
      if (selector.startsWith('xpath=')) {
        const xpath = selector.slice(6)
        const target = `::-p-xpath(${xpath})`
        await page.waitForSelector(target, { timeout: 2000 })
        await page.click(target, { clickCount: 3 })
        await page.keyboard.press('Backspace')
        await page.type(target, value)
        return true
      }
      await page.waitForSelector(selector, { timeout: 2000 })
      await page.click(selector, { clickCount: 3 })
      await page.keyboard.press('Backspace')
      await page.type(selector, value)
      return true
    } catch {}
  }
  return false
}

async function waitForManualLogin(page) {
  const start = Date.now()
  const timeoutMs = 10 * 60 * 1000

  while (Date.now() - start < timeoutMs) {
    const url = page.url()
    if (!url.includes('/sign') && !url.includes('/login') && !url.includes('/auth')) {
      return true
    }

    const text = await page.evaluate(() => document.body.innerText || '')
    if (/projects|create project|dashboard/i.test(text)) {
      return true
    }

    await new Promise(resolve => setTimeout(resolve, 2000))
  }

  throw new Error('Timed out waiting for manual Neon login')
}

async function maybeCreateProject(page) {
  await page.goto('https://console.neon.tech/app/projects', { waitUntil: 'networkidle2' })

  const createClicked = await clickFirst(page, [
    'button[data-testid="create-project-button"]',
    'xpath=//button[contains(., "Create project")]',
    'xpath=//a[contains(., "Create project")]',
    'xpath=//button[contains(., "New project")]',
  ])

  if (!createClicked) {
    return
  }

  await page.waitForTimeout(1500)

  await typeFirst(page, [
    'input[name="projectName"]',
    'input[id="projectName"]',
    'xpath=//input[contains(@placeholder, "Project")]',
    'xpath=//input[contains(@aria-label, "Project")]',
  ], 'vship')

  await clickFirst(page, [
    'xpath=//button[contains(., "Create project")]',
    'xpath=//button[contains(., "Create Project")]',
    'xpath=//button[contains(., "Continue")]',
  ])
}

async function tryCaptureConnectionString(page) {
  const start = Date.now()
  const timeoutMs = 3 * 60 * 1000

  while (Date.now() - start < timeoutMs) {
    const text = await page.evaluate(() => document.body.innerText || '')
    const match = text.match(/postgres(?:ql)?:\/\/[^\s"']+/i)
    if (match) {
      return match[0]
    }

    await clickFirst(page, [
      'xpath=//button[contains(., "Connection details")]',
      'xpath=//button[contains(., "Connection string")]',
      'xpath=//button[contains(., "Show")]',
      'xpath=//button[contains(., "Reveal")]',
      'xpath=//button[contains(., "Copy")]',
      'xpath=//a[contains(., "Connection details")]',
    ])

    await page.waitForTimeout(2000)
  }

  return null
}

async function main() {
  ensureDir(outputDir)

  const browser = await puppeteer.launch({
    executablePath: chromePath,
    headless: false,
    userDataDir,
    defaultViewport: { width: 1440, height: 960 },
    args: ['--start-maximized'],
  })

  const page = await browser.newPage()
  page.setDefaultTimeout(15000)

  console.log('Opening Neon login page...')
  await page.goto('https://console.neon.tech', { waitUntil: 'networkidle2' })

  console.log('Complete Neon login in the opened browser if needed...')
  await waitForManualLogin(page)

  console.log('Logged in. Trying to create/open a project...')
  await maybeCreateProject(page)

  console.log('Trying to find a PostgreSQL connection string...')
  const connectionString = await tryCaptureConnectionString(page)

  const result = {
    capturedAt: new Date().toISOString(),
    currentUrl: page.url(),
    connectionString,
  }

  fs.writeFileSync(outputFile, JSON.stringify(result, null, 2))
  console.log(`Saved result to ${outputFile}`)

  if (connectionString) {
    console.log(`Connection string: ${connectionString}`)
  } else {
    console.log('Connection string not detected automatically. Leave the browser open and retrieve it manually if needed.')
  }
}

main().catch(err => {
  console.error(err)
  process.exit(1)
})
