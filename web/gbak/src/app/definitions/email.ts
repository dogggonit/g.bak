export class GormModel {
  ID: number = 0;
  CreatedAt: Date = new Date();
  UpdatedAt: Date = new Date();
  DeletedAt: null | Date = null;

  constructor(gm: GormModel) {
    this.ID = gm.ID;
    this.CreatedAt = gm.CreatedAt;
    this.UpdatedAt = gm.UpdatedAt;
    this.DeletedAt = gm.DeletedAt;
  }
}

export class Email extends GormModel {
  readonly GmailId: string;
  readonly GmailThreadId: string;
  readonly Snippet: string;
  readonly GmailLabels: Label[] = [];
  readonly GmailDate: number;
  readonly EstimatedSize: number;
  readonly Message?: Message;

  constructor(e: Email) {
    super(e);
    this.GmailId = e.GmailId;
    this.GmailThreadId = e.GmailThreadId;
    this.Snippet = e.Snippet;
    this.GmailDate = e.GmailDate;
    this.EstimatedSize = e.EstimatedSize;
    this.Message = e.Message ? new Message(e.Message) : undefined;
    (e.GmailLabels || []).filter(l => !!l).forEach(l => this.GmailLabels.push(new Label(l)));
  }
}

export class Message extends GormModel {
  GmailId: string;
  MimeType: string;
  Filename: string;
  Body?: MessageBody;
  Headers: MessageHeader[] = [];
  Parts: Message[] = [];
  MessageID: number;
  EmailID: number;

  constructor(m: Message) {
    super(m);
    this.GmailId = m.GmailId;
    this.MimeType = m.MimeType;
    this.Filename = m.Filename;
    this.MessageID = m.MessageID;
    this.EmailID = m.EmailID;
    this.Body = m.Body ? new MessageBody(m.Body) : undefined;
    (m.Headers || []).filter(h => !!h).forEach(h => this.Headers.push(new MessageHeader(h)));
    (m.Parts || []).filter(p => !!p).forEach(p => this.Parts.push(new Message(p)));
  }

  download() {
    const data = this.Body?.Attachment?.Data;
    if (!data) {
      return;
    }

    const file = new window.Blob([data], { type: this.MimeType });
    const downloadAnchor = document.createElement("a");
    downloadAnchor.style.display = "none";
    downloadAnchor.href = URL.createObjectURL(file);
    downloadAnchor.download = this.Filename;
    downloadAnchor.click();
  }
}

export class MessageBody extends GormModel {
  MessageID: number;
  GmailAttachmentId: string;
  Data: string;
  Size: number;
  Attachment?: Attachment;

  constructor(mb: MessageBody) {
    super(mb);
    this.MessageID = mb.MessageID;
    this.GmailAttachmentId = mb.GmailAttachmentId;
    this.Data = mb.Data;
    this.Size = mb.Size;
    if (mb.Attachment) {
      this.Attachment = new Attachment(mb.Attachment);
    }
  }
}

export class MessageHeader extends GormModel {
  Name: string;
  Value: string;
  MessageID: number;

  constructor(mh: MessageHeader) {
    super(mh);
    this.Name = mh.Name;
    this.Value = mh.Value;
    this.MessageID = mh.MessageID;
  }
}

export class Label extends GormModel {
  GmailId: string;
  Label: string;
  LabelType: string;
  TextColor: string;
  BackgroundColor: string;

  constructor(l: Label) {
    super(l);
    this.GmailId = l.GmailId;
    this.Label = l.Label;
    this.LabelType = l.LabelType;
    this.TextColor = l.TextColor;
    this.BackgroundColor = l.BackgroundColor;
  }
}

export class Attachment extends GormModel {
  Filename: string;
  Data: Uint8Array;

  constructor(a: Attachment) {
    super(a);
    this.Filename = a.Filename;

    if (typeof a.Data as unknown === 'string') {
      const binary = atob(a.Data as unknown as string);
      const buffer = new ArrayBuffer(binary.length);
      this.Data = new Uint8Array(buffer);

      for (let i = 0; i < binary.length; i++) {
        this.Data[i] = binary.charCodeAt(i);
      }
    } else {
      this.Data = a.Data
    }
  }

  get base64(): string {
    let decoder = new TextDecoder('utf8');
    return btoa(decoder.decode(this.Data));
  }
}
