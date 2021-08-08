const puppeteer = require('puppeteer');
const repositoriesList = []

puppeteer.launch({ headless: true }).then(async browser => {
  const page = await browser.newPage()

  await page.setRequestInterception(true);
  page.on("request", requset => {
    if(requset.resourceType() === "document") {
      requset.continue()
    } else {
      requset.abort()
    }
  })

  await page.goto('https://github.com/trending/go?since=daily')

  const repositoriesName = await page.$$eval(".Box-row > .lh-condensed > a", item => item.map(e => e.text))
  const repositoriesOverview = await page.$$eval(".Box-row > p", item => item.map(e => e.innerHTML))

  await browser.close()

  for(let i = 0; i < 24; i++) {
    repositoriesList.push({
      "name": repositoriesName[i].replace(/[ ]|[\r\n]/g,""),
      "overview": repositoriesOverview[i].replace(/<[^>]+>|[\r\n]/g,"").trim()
    })
  }

  console.log(repositoriesList)
})