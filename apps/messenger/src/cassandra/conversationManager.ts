import { Client, types } from 'cassandra-driver';
import { Injectable } from '@nestjs/common';

@Injectable()
export class ConversationManager {
  private client: Client;

  constructor() {
    this.client = new Client({
      contactPoints: ['172.31.0.3'], // Adjust based on your Cassandra setup
      localDataCenter: 'my-datacenter-2', // Adjust based on your Cassandra setup
    });
  }

  async createSchema(schema: string): Promise<void> {
    try {
      await this.client.execute(
        `CREATE KEYSPACE IF NOT EXISTS ${schema} WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }`,
      );
      await this.client.execute(
        `CREATE TABLE IF NOT EXISTS ${schema}.conversations (id uuid PRIMARY KEY, title text, type int, members set<uuid>, image blob, created_at timestamp)`,
      );
      await this.client.execute(
        `CREATE INDEX IF NOT EXISTS ON ${schema}.conversations (members)`,
      );
    } catch (error) {
      throw new Error(`Failed to create schema: ${error.message}`);
    }
  }

  async getByUserID(
    userID: string,
    limit: number,
    offset: number,
  ): Promise<any[]> {
    try {
      const query =
        'SELECT id, title, type, members, image, created_at FROM conversations WHERE members CONTAINS ?';
      const result = await this.client.execute(query, [userID], {
        prepare: true,
      });
      return result.rows;
    } catch (error) {
      throw new Error(
        `Failed to get conversations by user ID: ${error.message}`,
      );
    }
  }

  async getByID(id: string): Promise<any> {
    try {
      const query =
        'SELECT id, title, type, members, image, created_at FROM conversations WHERE id = ? LIMIT 1';
      const result = await this.client.execute(
        query,
        [types.Uuid.fromString(id)],
        { prepare: true },
      );
      return result.first();
    } catch (error) {
      throw new Error(`Failed to get conversation by ID: ${error.message}`);
    }
  }

  async create(conversation: any): Promise<void> {
    try {
      const id = types.Uuid.random(); // Generate a new UUID for the conversation
      const query =
        'INSERT INTO conversations (id, title, type, members, image, created_at) VALUES (?, ?, ?, ?, ?, ?)';
      await this.client.execute(
        query,
        [
          id,
          conversation.title,
          conversation.type,
          conversation.members,
          conversation.image,
          conversation.created_at,
        ],
        { prepare: true },
      );
    } catch (error) {
      throw new Error(`Failed to create conversation: ${error.message}`);
    }
  }

  async removeMember(cID: string, memberID: string): Promise<void> {
    try {
      await this.client.execute(
        'UPDATE conversations SET members = members - { ? } WHERE id = ?',
        [memberID, cID],
        { prepare: true },
      );
    } catch (error) {
      throw new Error(`Failed to remove member: ${error.message}`);
    }
  }

  async addMember(cID: string, memberID: string): Promise<void> {
    try {
      await this.client.execute(
        'UPDATE conversations SET members = members + { ? } WHERE id = ?',
        [memberID, cID],
        { prepare: true },
      );
    } catch (error) {
      throw new Error(`Failed to add member: ${error.message}`);
    }
  }

  async update(conversation: any): Promise<void> {
    try {
      await this.client.execute(
        'UPDATE conversations SET title = ?, image = ? WHERE id = ?',
        [conversation.title, conversation.image, conversation.id],
        { prepare: true },
      );
    } catch (error) {
      throw new Error(`Failed to update conversation: ${error.message}`);
    }
  }
}
