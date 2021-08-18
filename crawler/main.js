const puppeteer = require("puppeteer")
const fetch = require('node-fetch')

const url = "http://192.168.0.21:8080/api/v1/trending/daily/save"

puppeteer.launch({ headless: true }).then(async browser => {
  try {
    const page = await browser.newPage()

    await page.setRequestInterception(true)
    page.on("request", requset => {
      if (requset.resourceType() === "document") {
        requset.continue()
      } else {
        requset.abort()
      }
    })

    let result = await query(page, "daily")

    await browser.close()

    post(result)
  } catch (err) {
    console.log(err)
  }
})

async function post(data) {
  try {
    const response = await fetch(url, {
      body: JSON.stringify(data),
      method: 'POST',
    })
    console.log(response.statusText)
  } catch (err) {
    return err
  }
}

async function query(page, dateRange) {
  let origin = {}
  let temp = []

  try {
    await page.goto("https://github.com/trending/go?since=" + dateRange, { waitUntil: "load" });

    [origin.nameAndUrl, origin.overview, origin.star, origin.fork, origin.todayStar, origin.boxCount] = await Promise.all([
      page.$$eval(".Box-row > h1 > a", item => item.map(e => { return { name: e.text, href: e.href } })),
      page.$$eval(".Box-row > p", item => item.map(e => e.innerHTML)),
      page.$$eval(".Box-row > .f6 > span + a", item => item.map(e => e.text)),
      page.$$eval(".Box-row > .f6 > span + a + a", item => item.map(e => e.text)),
      page.$$eval(".Box-row > .f6 > .float-sm-right", item => item.map(e => e.innerHTML)),
      page.$$eval(".Box-row", item => item.length),
    ])
  } catch (err) {
    console.log(err)
  }

  for (let i = 0; i < origin.boxCount; i++) {
    temp.push({
      "name": origin.nameAndUrl[i].name.replace(/[ ]|[\r\n]/g, ""),
      "url": origin.nameAndUrl[i].href.replace(/[ ]|[\r\n]/g, ""),
      "overview": origin.overview[i].replace(/<[^>]+>|[\r\n]/g, "").trim(),
      "star": +origin.star[i].replace(/[^\d.]/g, ""),
      "todayStar": +origin.todayStar[i].replace(/<[^>]+>|[\r\n]|[^\d.]/g, ""),
      "fork": +origin.fork[i].replace(/[^\d.]/g, ""),
    })
  }

  return {
    date: new Date(),
    [dateRange]: temp
  }
}