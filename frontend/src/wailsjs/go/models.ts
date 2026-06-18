export namespace speedtest_util {
	
	export class HistoryEntry {
	    // Go type: time
	    timestamp: any;
	    server: string;
	    ping: number;
	    download: number;
	    upload: number;
	
	    static createFrom(source: any = {}) {
	        return new HistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.server = source["server"];
	        this.ping = source["ping"];
	        this.download = source["download"];
	        this.upload = source["upload"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace updater {

	export class UpdateInfo {
	    LatestVersion: string;
	    ReleasePageURL: string;
	    AssetSizeBytes: number;
	    HasUpdate: boolean;
	    DownloadURL: string;

	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LatestVersion = source["LatestVersion"];
	        this.ReleasePageURL = source["ReleasePageURL"];
	        this.AssetSizeBytes = source["AssetSizeBytes"];
	        this.HasUpdate = source["HasUpdate"];
	        this.DownloadURL = source["DownloadURL"];
	    }
	}

}
