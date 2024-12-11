export namespace api {
	
	export class ApplyForResponse {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    UserID: number;
	    FriendID: number;
	    IsGroup: boolean;
	    IsConsent: number;
	    applicant_id: number;
	    applicant_username: string;
	    applicant_avatar: string;
	    friend_id: number;
	    friend_username: string;
	    friend_avatar: string;
	    friend_uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new ApplyForResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.UserID = source["UserID"];
	        this.FriendID = source["FriendID"];
	        this.IsGroup = source["IsGroup"];
	        this.IsConsent = source["IsConsent"];
	        this.applicant_id = source["applicant_id"];
	        this.applicant_username = source["applicant_username"];
	        this.applicant_avatar = source["applicant_avatar"];
	        this.friend_id = source["friend_id"];
	        this.friend_username = source["friend_username"];
	        this.friend_avatar = source["friend_avatar"];
	        this.friend_uuid = source["friend_uuid"];
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
	export class Info {
	    username: string;
	    password: string;
	    signature: string;
	    avatar: string;
	    uuid: string;
	    id: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	        this.signature = source["signature"];
	        this.avatar = source["avatar"];
	        this.uuid = source["uuid"];
	        this.id = source["id"];
	    }
	}
	export class FriendRes {
	    InfoList: Info[];
	    Err: any;
	
	    static createFrom(source: any = {}) {
	        return new FriendRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.InfoList = this.convertValues(source["InfoList"], Info);
	        this.Err = source["Err"];
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
	export class GroupMessageList {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    content: string;
	    send_user_id: number;
	    receiver_user_id: number;
	    room_id: string;
	    file_url: string;
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    Username: string;
	    Password: string;
	    Signature: string;
	    Avatar: string;
	    UUID: string;
	
	    static createFrom(source: any = {}) {
	        return new GroupMessageList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.content = source["content"];
	        this.send_user_id = source["send_user_id"];
	        this.receiver_user_id = source["receiver_user_id"];
	        this.room_id = source["room_id"];
	        this.file_url = source["file_url"];
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.Signature = source["Signature"];
	        this.Avatar = source["Avatar"];
	        this.UUID = source["UUID"];
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
	
	export class LoginResponse {
	    message: string;
	    userID: number;
	    UUID: string;
	    useravatar: string;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new LoginResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.message = source["message"];
	        this.userID = source["userID"];
	        this.UUID = source["UUID"];
	        this.useravatar = source["useravatar"];
	        this.status = source["status"];
	    }
	}
	export class UserInfo {
	    userName: string;
	    userID: number;
	    avatar: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.userName = source["userName"];
	        this.userID = source["userID"];
	        this.avatar = source["avatar"];
	    }
	}
	export class Room {
	    roomID: string;
	    roomName: string;
	    userInfo: UserInfo[];
	    createUser: number;
	
	    static createFrom(source: any = {}) {
	        return new Room(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.roomID = source["roomID"];
	        this.roomName = source["roomName"];
	        this.userInfo = this.convertValues(source["userInfo"], UserInfo);
	        this.createUser = source["createUser"];
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
	export class Response {
	    room: Room[];
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.room = this.convertValues(source["room"], Room);
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
	
	
	export class UserSearchResult {
	    ID: number;
	    // Go type: time
	    CreatedAt: any;
	    // Go type: time
	    UpdatedAt: any;
	    // Go type: gorm
	    DeletedAt: any;
	    Username: string;
	    Password: string;
	    Signature: string;
	    Avatar: string;
	    UUID: string;
	    is_friend: number;
	
	    static createFrom(source: any = {}) {
	        return new UserSearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.CreatedAt = this.convertValues(source["CreatedAt"], null);
	        this.UpdatedAt = this.convertValues(source["UpdatedAt"], null);
	        this.DeletedAt = this.convertValues(source["DeletedAt"], null);
	        this.Username = source["Username"];
	        this.Password = source["Password"];
	        this.Signature = source["Signature"];
	        this.Avatar = source["Avatar"];
	        this.UUID = source["UUID"];
	        this.is_friend = source["is_friend"];
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

