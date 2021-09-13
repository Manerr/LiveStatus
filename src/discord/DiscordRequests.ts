import {ConfigModule} from "../model/Config.ts";
import {Requests} from "../model/Requests.ts";
import HttpMethod = Requests.HttpMethod;
import DiscordConfig = ConfigModule.DiscordConfig;

export default class DiscordRequests {

    private readonly discordConfig: DiscordConfig;

    public constructor(discordConfig: DiscordConfig) {
        this.discordConfig = discordConfig;
    }

    public async createMessage(channelId: string, body: any) {
        return await this.fetchDiscord(`channels/${channelId}/messages`, HttpMethod.POST, body).then(res => res.json());
    }

    public async editMessage(channelId: string, messageId: string, body: any) {
        return await this.fetchDiscord(`channels/${channelId}/messages/${messageId}`, HttpMethod.PATCH, body).then(res => res.json());
    }

    private async fetchDiscord(apiPath: string, httpMethod: HttpMethod, body: any) {
        const headers: Headers = new Headers({
            'Authorization': 'Bot ' + this.discordConfig.discordToken,
            'Content-Type': 'application/json'
        });
        return await Requests.fetchURL(`https://discord.com/api/${apiPath}`, httpMethod, headers, body);
    }
}