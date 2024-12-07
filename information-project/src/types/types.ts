export interface InputSearch {
  page: number
  word: string
}

export const NUMBER_ROWS=10

export interface IEmailInformation{
  _source:IEmail
}

export interface IEmail {
  MessageID: string
  Date: string
  From: string
  To: string
  Subject: string
  Mime_Version: string
  Content_Type: string
  Content_Transfer_Encoding: string
  X_From: string
  X_To: string
  X_cc: string
  X_bcc: string
  X_Folder: string
  X_Origin: string
  X_FileName: string
  Body: string
}

export type ISearchRequest = {
  query: {
    term: string;
  };
  from: number;
  size: number;
  sort_fields: string;
};

export type IResponseInformation = {
  hits: {
    total: {
      value: number
    }
    hits: [IEmailInformation]
  }
}