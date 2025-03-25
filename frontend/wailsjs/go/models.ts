export namespace main {
	
	export class FileInfo {
	    name: string;
	    path: string;
	    pages: number;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.pages = source["pages"];
	    }
	}
	export class ProcessResult {
	    file: string;
	    success: boolean;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.success = source["success"];
	        this.message = source["message"];
	    }
	}

}

