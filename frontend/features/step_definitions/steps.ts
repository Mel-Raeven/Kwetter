import * as assert from "assert";
import { Given, When, Then } from "@cucumber/cucumber";
import { postMessage } from "../../src/functions/postMessage.js";

interface Token {
  idToken: string;
}

interface Response {
  PostMessageResponse: PostMessageResponse;
}

interface PostMessageResponse {
  success: boolean;
  data?: any;
  error?: string;
}

Given("the user is logged in with a valid ID token", function(this: Token) {
  this.idToken = "eyJraWQiOiJMUlZIWjJ6SU9jSkwrSDJmcHhmREg0UHJpbXVJQWdHVUgyM2hqMWpOVVRFPSIsImFsZyI6IlJTMjU2In0.eyJhdF9oYXNoIjoidUpXV2tUa2pmaXYxVTl3STZEMlhfQSIsInN1YiI6IjEzNjQyOGYyLTgwNTEtNzA5OS1kN2EyLTQ2NWFjYzE1OTNiYyIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJpc3MiOiJodHRwczpcL1wvY29nbml0by1pZHAuZXUtY2VudHJhbC0xLmFtYXpvbmF3cy5jb21cL2V1LWNlbnRyYWwtMV9qeERPSTBPOEQiLCJjb2duaXRvOnVzZXJuYW1lIjoiMTM2NDI4ZjItODA1MS03MDk5LWQ3YTItNDY1YWNjMTU5M2JjIiwib3JpZ2luX2p0aSI6IjNmYjU5Yzg1LTNkZjItNDcxNC05OGQ3LTk0ZWY1NjVmN2ZmMCIsImF1ZCI6IjF2bmxkNXZsN3RvY21lY3Zia2prc21zZjJyIiwiZXZlbnRfaWQiOiI5MzM3OWNiOC05Mzk4LTQzZDEtYmFiZS0wZTRjZjYxNmE5N2QiLCJ0b2tlbl91c2UiOiJpZCIsImF1dGhfdGltZSI6MTcxNzAxNzQ0MSwiZXhwIjoxNzE3MDIxMDQxLCJpYXQiOjE3MTcwMTc0NDEsImp0aSI6IjMwM2RkNDJlLWIzMDMtNGI1OC04Mzg1LTc1ODFiNzUyMTIyOCIsImVtYWlsIjoibWVsLnJhZXZlbkBnbWFpbC5jb20ifQ.GR5WdOI4tqO35uVYkC3Is0qhwEsZE2CcVGXHZAv1KIPkY3c1tSWt6Q2IQJTIFIJk3LlTdSnq6GzRsYeTLmkS2sFLRrtsVfY305VHRebnCock5bSHiwJP93ulJazTyJI4kiERirGHqMhigZnZzb60nvVx1fyZKtnXxrlnin1Cy75It5gC_mr3oJNDPLB00xwkDbg8KQG_5UlIuHRedq5zT18akwMmZa-s2WpDvqeSQ-lNFXoNSWNZ-e2SMI-k9G9EurrPpZYP-rzN3jgeOM3tQx1ikWObOtsjPXafhJwJWa9RikWSHTmmWuzQ3MbCdWDWA1bl-jWOLSB1z1O9xMvMHg";
});

When("the user posts {string}", async function(this: Token & Response, message: string) {
  try {
    this.PostMessageResponse = await postMessage(this.idToken, message);
  } catch (error) {
    console.error("Error posting message:", error);
    throw error;
  }
});

Then("the message should be sent", function(this: Response) {
  assert.equal(this.PostMessageResponse.success, true);
});
