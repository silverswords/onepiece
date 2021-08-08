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

  const nameAndUrl = await page.$$eval(".Box-row > h1 > a", item => item.map(e => {return {name: e.text, href: e.href}}))
  const overview = await page.$$eval(".Box-row > p", item => item.map(e => e.innerHTML))
  const cumulativeStar = await page.$$eval(".Box-row > .f6 > span + a", item => item.map(e => e.text))
  const fork = await page.$$eval(".Box-row > .f6 > span + a + a", item => item.map(e => e.text))
  const star = await page.$$eval(".Box-row > .f6 > .float-sm-right", item => item.map(e => e.innerHTML))

  await browser.close()

  for(let i = 0; i < 25; i++) {
    repositoriesList.push({
      "name": nameAndUrl[i].name.replace(/[ ]|[\r\n]/g,""),
      "url": nameAndUrl[i].href.replace(/[ ]|[\r\n]/g,""),
      "overview": overview[i].replace(/<[^>]+>|[\r\n]/g,"").trim(),
      "cumulativeStar": +cumulativeStar[i].replace(/[^\d.]/g,""),
      "star": +star[i].replace(/<[^>]+>|[\r\n]|[^\d.]/g,""),
      "fork": +fork[i].replace(/[^\d.]/g,""),
    })
  }

  console.log(repositoriesList)
})