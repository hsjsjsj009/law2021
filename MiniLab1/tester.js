'use strict'

const autocannon = require("autocannon")

autocannon({
    url:"http://localhost:8000/susunjadwal/api/scrap-schedule",
    headers:{
        "Authorization":"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiNjBkZDhmMTVkN2EwMTgyYzYyN2QxYTNjIiwibWFqb3JfaWQiOiI2MGRkOGMwZTdiNjYxOGY1ZTQ2YmYwYmMifQ.MBR5tU7HfKGLS2r-y-4CzdS2wYs0Q69pwIha2qyyFeI",
        "Content-Type":"application/json"
    },
    body:JSON.stringify({
        'username':'dipta.laksmana',
        'password':'baswara546456'
    }),
    method:"POST",
    connections:10,
    pipelining:10,
    duration:2
},console.log)