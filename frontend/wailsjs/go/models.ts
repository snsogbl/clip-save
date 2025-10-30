export namespace common {
	
	export class ClipboardItem {
	    ID: string;
	    Content: string;
	    ContentType: string;
	    ContentHash: string;
	    ImageData: number[];
	    FilePaths: string;
	    FileInfo: string;
	    // Go type: time
	    Timestamp: any;
	    Source: string;
	    CharCount: number;
	    WordCount: number;
	    IsFavorite: number;
	
	    static createFrom(source: any = {}) {
	        return new ClipboardItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Content = source["Content"];
	        this.ContentType = source["ContentType"];
	        this.ContentHash = source["ContentHash"];
	        this.ImageData = source["ImageData"];
	        this.FilePaths = source["FilePaths"];
	        this.FileInfo = source["FileInfo"];
	        this.Timestamp = this.convertValues(source["Timestamp"], null);
	        this.Source = source["Source"];
	        this.CharCount = source["CharCount"];
	        this.WordCount = source["WordCount"];
	        this.IsFavorite = source["IsFavorite"];
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
	export class FileInfo {
	    name: string;
	    path: string;
	    size: number;
	    is_dir: boolean;
	    exists: boolean;
	    extension: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.is_dir = source["is_dir"];
	        this.exists = source["exists"];
	        this.extension = source["extension"];
	    }
	}

}

