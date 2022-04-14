export class APIEmails {
  Total: number;
  Emails: APIEmail[] = [];

  constructor(ae: APIEmails) {
    this.Total = ae.Total;
    (ae.Emails || []).forEach(e => this.Emails.push(new APIEmail(e)));
  }
}

export class APIEmail {
  ID: number;
  Snippet: string;
  To: string;
  From: string;
  Subject: string;
  Date: number;
  HasAttachment: boolean;
  Labels: APILabel[] = [];

  constructor(e: APIEmail) {
    this.ID = e.ID;
    this.Snippet = e.Snippet
    this.To = e.To;
    this.From = e.From;
    this.Subject = e.Subject;
    this.Date = e.Date;
    this.HasAttachment = e.HasAttachment;
    (e.Labels || []).forEach(l => this.Labels.push(new APILabel(l)));
  }
}

export class APIUser {
  Email: string;
  Labels: APILabel[] = [];

  constructor(u: APIUser) {
    this.Email = u.Email;
    (u.Labels || []).forEach(l => this.Labels.push(new APILabel(l)))
  }
}

export class APILabel {
  GmailId: string;
  Label: string;
  TextColor: string;
  BackgroundColor: string;

  constructor(l: APILabel) {
    this.GmailId = l.GmailId;
    this.Label = l.Label;
    this.TextColor = l.TextColor;
    this.BackgroundColor = l.BackgroundColor;
  }
}
