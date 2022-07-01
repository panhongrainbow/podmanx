"use strict";
import {proj} from "../proj";

import http from "http";
import express from "express";


export async function start() {
    const app = express();
    const server = http.createServer(app);

    // Routing
    app.use(express.static(proj.config.basedir + "/public"));
    app.get("/healthz", function(req, res) {
        res.send("ok@"+Date.now());
    });

    server.listen(proj.config.LISTEN_PORT, proj.config.LISTEN_HOST, function() {
        console.warn(`listening at port ${proj.config.LISTEN_PORT}`);
    });
}
