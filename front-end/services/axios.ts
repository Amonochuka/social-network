import loadEnvFile from "@/config/config";
import axios from "axios";

export const config = loadEnvFile({ urlType: "baseUrl" });

export const Api = axios.create({
    baseURL: config.baseUrl,
    withCredentials: true, // sends the session_id cookie automatically
    headers: {
        "Content-Type": "application/json",
    },
});