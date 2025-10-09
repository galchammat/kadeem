export namespace models {
	
	export class LeagueOfLegendsAccount {
	    puuid: string;
	    tagLine: string;
	    gameName: string;
	    region?: string;
	    streamer?: string;
	
	    static createFrom(source: any = {}) {
	        return new LeagueOfLegendsAccount(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.puuid = source["puuid"];
	        this.tagLine = source["tagLine"];
	        this.gameName = source["gameName"];
	        this.region = source["region"];
	        this.streamer = source["streamer"];
	    }
	}

}

