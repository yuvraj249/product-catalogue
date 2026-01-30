import styled from "styled-components";

export const PageHeader = styled.div`
  margin-bottom: 32px;
`;

export const Title = styled.h1`
  font-size: 32px;
  font-weight: 700;
  color: #1f2d36;
  margin: 0;
`;

export const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
`;

export const StatCard = styled.div`
  background: linear-gradient(135deg, #ffffff, #f3f8fb);
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  padding: 24px;

  display: flex;
  align-items: center;
  gap: 16px;

  cursor: pointer;
  transition: all 0.25s ease;

  &:hover {
    border-color: rgba(87, 199, 133, 0.6);
    box-shadow: 0 10px 25px rgba(42, 123, 155, 0.18);
    transform: translateY(-2px);
  }
`;

export const StatIcon = styled.div`
  width: 56px;
  height: 56px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: linear-gradient(
    135deg,
    rgba(42, 123, 155, 0.18),
    rgba(87, 199, 133, 0.18)
  );

  border-radius: 8px;
`;

export const StatInfo = styled.div`
  flex: 1;
`;

export const StatValue = styled.div`
  font-size: 24px;
  font-weight: 700;
  color: #1f2d36;
  margin-bottom: 4px;
`;

export const StatLabel = styled.div`
  font-size: 14px;
  color: #4f6b72;
`;

export const Section = styled.div`
  background: linear-gradient(135deg, #ffffff, #f3f8fb);
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  padding: 32px;

  @media (max-width: 768px) {
    padding: 20px;
  }
`;

export const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
`;

export const SectionTitle = styled.h2`
  font-size: 18px;
  font-weight: 600;
  color: #1f2d36;
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
`;

export const ThresholdBadge = styled.div`
  padding: 6px 12px;
  background: rgba(87, 199, 133, 0.15);
  border: 1px solid rgba(87, 199, 133, 0.3);
  border-radius: 6px;
  color: #2a7b9b;
  font-size: 12px;
  font-weight: 500;
`;

export const TableWrapper = styled.div`
  background: #ffffff;
  border: 1px solid rgba(42, 123, 155, 0.25);
  border-radius: 12px;
  overflow: hidden;
  margin-top: 24px;
`;

export const Table = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

export const Th = styled.th`
  padding: 12px 16px;
  background: #f3f8fb;
  color: #4f6b72;
  font-size: 14px;
  font-weight: 600;
  text-align: left;
  border-bottom: 2px solid rgba(42, 123, 155, 0.25);
`;

export const Td = styled.td`
  padding: 16px;
  border-bottom: 1px solid rgba(42, 123, 155, 0.18);
  color: #1f2d36;
  font-size: 14px;
  font-weight: 500;
`;

export const Icon = styled.img`
  width: 20px;
  height: 20px;
`;
