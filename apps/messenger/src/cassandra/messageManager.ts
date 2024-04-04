import { Client, types } from 'cassandra-driver';
import { Injectable } from '@nestjs/common';

@Injectable()
export class MessageManager {
  private client: Client;

  constructor() {
    this.client = new Client({
      contactPoints: ['172.31.0.4', '172.31.0.2'], // Adjust based on your Cassandra setup
      localDataCenter: 'my-datacenter-1', // Adjust based on your Cassandra setup
    });
  }

  async createSchema(schema: string): Promise<void> {
    try {
      await this.client.execute(
        `CREATE KEYSPACE IF NOT EXISTS ${schema} WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 2 }`,
      );
      await this.client.execute(
        `CREATE TABLE IF NOT EXISTS ${schema}.messages (
          id uuid,
          sender_id uuid,
          conversation_id uuid,
          text text,
          audio blob,
          image blob,
          event int,
          extra_params map<text, text>,
          readed map<uuid, timestamp>,
          created_at timestamp,
          delivery_at timestamp,
          PRIMARY KEY ((id), conversation_id)
        )`,
      );
      await this.client.execute(
        `CREATE INDEX IF NOT EXISTS ON ${schema}.messages (conversation_id)`,
      );
    } catch (error) {
      throw new Error(`Failed to create schema: ${error.message}`);
    }
  }

  async getMessages(
    conversationID: string,
    limit: number,
    offset: number,
  ): Promise<any[]> {
    try {
      const query = `SELECT id, sender_id, text, audio, image, event, extra_params, created_at, delivery_at FROM messages_new.messages WHERE conversation_id = ? LIMIT ?`;
      const result = await this.client.execute(query, [conversationID, limit], {
        prepare: true,
      });

      return result.rows.map((row) => ({
        ...row,
        id: row.id.toString(),
        senderId: row.sender_id.toString(),
        createdAt: row.created_at,
        deliveryAt: row.delivery_at,
      }));
    } catch (error) {
      throw new Error(`Failed to get messages: ${error.message}`);
    }
  }

  async saveMessage(message: any, conversationID: string): Promise<void> {
    try {
      const id = types.Uuid.random();
      const query = `INSERT INTO messages_new.messages (id, sender_id, conversation_id, text, audio, image, event, extra_params, readed, created_at, delivery_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`;
      await this.client.execute(
        query,
        [
          id,
          message.senderId,
          conversationID,
          message.text,
          message.audio,
          message.image,
          message.event,
          message.extraParams,
          null, // Assuming 'readed' is not provided
          1712226768,
          1712226796,
        ],
        { prepare: true },
      );
    } catch (error) {
      throw new Error(`Failed to save message: ${error.message}`);
    }
  }
  async readMessages(conversationID: string, userID: string): Promise<void> {
    throw new Error('Not implemented');
  }
}
