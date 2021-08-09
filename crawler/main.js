const puppeteer = require("puppeteer");
const query = require("./query")
const fetch = require('node-fetch')
const url = "http://192.168.0.21:8080/api/v1/trending/daily/save"

puppeteer.launch({ headless: true }).then(async browser => {
  const page = await browser.newPage()

  await page.setRequestInterception(true);
  page.on("request", requset => {
    if (requset.resourceType() === "document") {
      requset.continue()
    } else {
      requset.abort()
    }
  })

  let result = await query.query(page, "daily")

  await browser.close()

  post(result)
})

async function post(data) {
  try {
    const response = await fetch(url, {
      body: JSON.stringify(data),
      method: 'POST',
    })
    console.log(response.statusText)
  } catch (err) {
    return err
  }
}