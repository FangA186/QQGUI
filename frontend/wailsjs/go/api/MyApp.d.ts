// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {api} from '../models';
import {model} from '../models';

export function ConsentFriend(arg1:string,arg2:string,arg3:string):Promise<string>;

export function CrateRoomReceive(arg1:string,arg2:Array<{[key: string]: string}>,arg3:string,arg4:string):Promise<void>;

export function CreateRoom(arg1:{[key: string]: {[key: string]: string}}):Promise<{[key: string]: any}>;

export function FriendList(arg1:string):Promise<api.FriendRes>;

export function GetApplicationRecord(arg1:string):Promise<Array<api.ApplyForResponse>>;

export function GetConsentList(arg1:string):Promise<Array<api.ApplyForResponse>>;

export function GetIsSpeakUserInfo(arg1:string):Promise<Array<api.ApplyForResponse>>;

export function HandleAddFriendRequest(arg1:string,arg2:string,arg3:boolean):Promise<string>;

export function IsSpeak(arg1:string,arg2:string):Promise<{[key: string]: any}>;

export function Login(arg1:string,arg2:string):Promise<api.LoginResponse>;

export function MessageList(arg1:string):Promise<Array<model.Message>>;

export function QueryRoom(arg1:number):Promise<api.Response>;

export function Receive(arg1:string,arg2:string,arg3:string,arg4:string):Promise<void>;

export function RegisterFunc(arg1:string,arg2:string):Promise<{[key: string]: any}>;

export function SearchUser(arg1:string,arg2:string):Promise<Array<api.UserSearchResult>>;

export function StartHTTPServer():Promise<void>;
